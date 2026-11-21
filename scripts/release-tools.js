const fs = require('fs');
const path = require('path');
const { spawnSync } = require('child_process');

/*
 * Release decision tooling used by .github/workflows/release.yml.
 *
 * Subcommands:
 *   check-version <tag>             tag (without "v"), package.json and constants.go must agree
 *   channel <version>               prints the npm dist-tag for a version
 *   publish <package-dir> <channel> publishes a package unless that version already exists
 *   notes <tag>                     prints .github/release-notes/<tag>.md when it has content
 *
 * SECURITY NOTE: commands are run with 'child_process.spawnSync' and an argument
 * array. We NEVER use 'shell: true' (same rule as index.js).
 */

const repoRoot = path.join(__dirname, '..');

// Semantic Versioning 2.0.0 (https://semver.org), without a leading "v".
const semverRegex = /^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$/;

// Same declaration style as scripts/sync-version.js.
const versionDeclRegex = /(?:const|var)\s+Version\s+=\s+"([^"]+)"/;

function parseVersion(version) {
    const match = semverRegex.exec(version);
    if (!match) {
        throw new Error(`"${version}" is not a valid version`);
    }
    return { prerelease: match[4] || '' };
}

function checkVersion(tag, packageVersion, builtinVersion) {
    const tagVersion = tag.startsWith('v') ? tag.slice(1) : tag;
    if (!semverRegex.test(tagVersion)) {
        throw new Error(`tag "${tag}" is not a valid version`);
    }
    if (tagVersion !== packageVersion || tagVersion !== builtinVersion) {
        throw new Error(
            `version mismatch: tag=${tag} package.json=${packageVersion} built-in=${builtinVersion}`
        );
    }
}

// Mirrors ChannelOf in src/internal/upgrade/provider_npm.go: the first
// pre-release identifier, or "latest" for stable versions.
function channelOf(version) {
    const { prerelease } = parseVersion(version);
    if (prerelease === '') {
        return 'latest';
    }
    const channel = prerelease.split('.')[0];
    if (/^\d+$/.test(channel)) {
        throw new Error(
            `pre-release "${version}" needs a named pre-release channel (e.g. ${version.split('-')[0]}-rc.1)`
        );
    }
    return channel;
}

function firstLine(text) {
    const line = String(text || '').split('\n').map((l) => l.trim()).find((l) => l !== '');
    return line || 'no output';
}

function defaultRun(cmd, args, options) {
    return spawnSync(cmd, args, { encoding: 'utf8', ...options });
}

function isNotFound(stderr) {
    return /\bE404\b|404 Not Found|is not in this registry/i.test(stderr || '');
}

function publishPackage(packageDir, channel, deps = {}) {
    const run = deps.run || defaultRun;
    const readFile = deps.readFile || ((p) => fs.readFileSync(p, 'utf8'));
    const log = deps.log || ((line) => console.log(line));

    const pkg = JSON.parse(readFile(path.join(packageDir, 'package.json')));
    if (!pkg.name || !pkg.version) {
        throw new Error(`${path.join(packageDir, 'package.json')} has no name or version`);
    }
    const spec = `${pkg.name}@${pkg.version}`;

    const view = run('npm', ['view', spec, 'version'], { cwd: packageDir });
    if (view.error) {
        throw new Error(`cannot query ${spec}: ${firstLine(view.error.message)}`);
    }

    let exists;
    if (view.status === 0) {
        // `npm view` for a missing version of an existing package exits 0 with
        // empty output instead of E404, so empty output also means absent.
        exists = String(view.stdout || '').trim() !== '';
    } else if (isNotFound(view.stderr)) {
        exists = false;
    } else {
        throw new Error(`cannot query ${spec}: ${firstLine(view.stderr)}`);
    }

    if (exists) {
        log(`skipped ${spec} (already published)`);
        return;
    }

    const publish = run(
        'npm',
        ['publish', '--access', 'public', '--provenance', '--tag', channel],
        { cwd: packageDir, stdio: ['ignore', 'inherit', 'pipe'] }
    );
    if (publish.error) {
        throw new Error(`failed to publish ${spec}: ${firstLine(publish.error.message)}`);
    }
    if (publish.status !== 0) {
        throw new Error(`failed to publish ${spec}: ${firstLine(publish.stderr)}`);
    }
    log(`published ${spec} (${channel})`);
}

function notesPath(tag, deps = {}) {
    const root = deps.root || repoRoot;
    const readFile = deps.readFile || ((p) => fs.readFileSync(p, 'utf8'));
    const file = path.join(root, '.github', 'release-notes', `${tag}.md`);
    let content;
    try {
        content = readFile(file);
    } catch (error) {
        if (error.code === 'ENOENT') {
            return null;
        }
        throw error;
    }
    return content.trim() === '' ? null : file;
}

function readRepoVersions() {
    const packageJson = JSON.parse(fs.readFileSync(path.join(repoRoot, 'package.json'), 'utf8'));
    const constantsGo = fs.readFileSync(
        path.join(repoRoot, 'src', 'internal', 'core', 'constants.go'),
        'utf8'
    );
    const match = versionDeclRegex.exec(constantsGo);
    if (!match) {
        throw new Error('could not find Version declaration in constants.go');
    }
    return { packageVersion: packageJson.version, builtinVersion: match[1] };
}

function requireArgs(args, count, usage) {
    if (args.length !== count) {
        throw new Error(`usage: release-tools.js ${usage}`);
    }
}

function main(argv) {
    const [command, ...args] = argv;
    switch (command) {
        case 'check-version': {
            requireArgs(args, 1, 'check-version <tag>');
            const { packageVersion, builtinVersion } = readRepoVersions();
            checkVersion(args[0], packageVersion, builtinVersion);
            console.log(`version ok: ${args[0]}`);
            break;
        }
        case 'channel':
            requireArgs(args, 1, 'channel <version>');
            console.log(channelOf(args[0]));
            break;
        case 'publish':
            requireArgs(args, 2, 'publish <package-dir> <channel>');
            publishPackage(path.resolve(args[0]), args[1]);
            break;
        case 'notes': {
            requireArgs(args, 1, 'notes <tag>');
            const file = notesPath(args[0]);
            if (file) {
                console.log(file);
            }
            break;
        }
        default:
            throw new Error(
                `unknown command "${command || ''}" (expected check-version, channel, publish or notes)`
            );
    }
}

if (require.main === module) {
    try {
        main(process.argv.slice(2));
    } catch (error) {
        console.error(`Error: ${firstLine(error.message)}`);
        process.exit(1);
    }
}

module.exports = {
    checkVersion,
    channelOf,
    publishPackage,
    notesPath,
};
