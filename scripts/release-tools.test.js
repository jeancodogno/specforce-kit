const test = require('node:test');
const assert = require('node:assert');
const path = require('path');

const {
    checkVersion,
    channelOf,
    publishPackage,
    notesPath,
    NATIVE_PACKAGES,
    checkNativeVersions,
    checkRelease,
    syncNativeVersions,
} = require('./release-tools');

// --- Version agreement (US-2) ---

test('checkVersion passes when tag, package.json and built-in versions agree', () => {
    assert.doesNotThrow(() => checkVersion('v1.0.0', '1.0.0', '1.0.0'));
});

test('checkVersion fails on mismatch naming tag, package.json and built-in values', () => {
    assert.throws(
        () => checkVersion('v1.0.0', '1.0.0-rc.1', '1.0.0'),
        (err) => {
            assert.match(err.message, /v1\.0\.0/);
            assert.match(err.message, /1\.0\.0-rc\.1/);
            assert.match(err.message, /package\.json/);
            assert.match(err.message, /built-in/);
            assert.ok(!err.message.includes('\n'), 'message must be a single line');
            return true;
        }
    );
});

test('checkVersion fails when the built-in version differs', () => {
    assert.throws(() => checkVersion('v1.0.0', '1.0.0', '0.9.0'), /0\.9\.0/);
});

test('checkVersion fails when the tag is not a valid version', () => {
    assert.throws(() => checkVersion('vnext', '1.0.0', '1.0.0'), /vnext.*not a valid version/);
});

// --- Channel derivation (US-3) ---

test('channelOf returns latest for a stable version', () => {
    assert.strictEqual(channelOf('1.0.0'), 'latest');
});

test('channelOf returns the first pre-release identifier', () => {
    assert.strictEqual(channelOf('1.0.0-beta.3'), 'beta');
    assert.strictEqual(channelOf('1.1.0-next.1'), 'next');
    assert.strictEqual(channelOf('1.0.0-rc.1'), 'rc');
});

test('channelOf fails for a numeric first pre-release identifier', () => {
    assert.throws(() => channelOf('1.1.0-1'), /named pre-release channel/);
});

test('channelOf fails for an invalid version', () => {
    assert.throws(() => channelOf('next'), /not a valid version/);
});

// --- Publish or skip (US-4, US-5) ---

const PKG_DIR = path.join('fake', 'npm', 'linux-x64');
const PKG_JSON = JSON.stringify({ name: '@scope/pkg-linux-x64', version: '1.0.0-rc.1' });

function fakeReader(files) {
    return (p) => {
        if (Object.prototype.hasOwnProperty.call(files, p)) {
            return files[p];
        }
        const err = new Error(`ENOENT: no such file or directory, open '${p}'`);
        err.code = 'ENOENT';
        throw err;
    };
}

function fakeRunner(responses) {
    const calls = [];
    const run = (cmd, args, options) => {
        calls.push({ cmd, args, options });
        const key = args[0];
        const response = responses[key];
        if (!response) {
            throw new Error(`unexpected command: ${cmd} ${args.join(' ')}`);
        }
        return { status: 0, stdout: '', stderr: '', ...response };
    };
    return { run, calls };
}

function setup(responses) {
    const runner = fakeRunner(responses);
    const logs = [];
    const deps = {
        run: runner.run,
        readFile: fakeReader({ [path.join(PKG_DIR, 'package.json')]: PKG_JSON }),
        log: (line) => logs.push(line),
    };
    return { runner, logs, deps };
}

test('publishPackage publishes an absent package with provenance and channel tag', () => {
    const { runner, logs, deps } = setup({
        view: { status: 1, stderr: 'npm error code E404\nnpm error 404 Not Found' },
        publish: { status: 0 },
    });

    publishPackage(PKG_DIR, 'rc', deps);

    assert.deepStrictEqual(runner.calls[0].args, ['view', '@scope/pkg-linux-x64@1.0.0-rc.1', 'version']);
    const publishCall = runner.calls.find((c) => c.args[0] === 'publish');
    assert.ok(publishCall, 'npm publish must be called');
    assert.strictEqual(publishCall.cmd, 'npm');
    assert.deepStrictEqual(publishCall.args, ['publish', '--access', 'public', '--provenance', '--tag', 'rc']);
    assert.strictEqual(publishCall.options.cwd, PKG_DIR);
    assert.deepStrictEqual(logs, ['published @scope/pkg-linux-x64@1.0.0-rc.1 (rc)']);
});

test('publishPackage treats an empty successful query as absent', () => {
    const { runner, logs, deps } = setup({
        view: { status: 0, stdout: '\n' },
        publish: { status: 0 },
    });

    publishPackage(PKG_DIR, 'rc', deps);

    assert.ok(runner.calls.some((c) => c.args[0] === 'publish'));
    assert.deepStrictEqual(logs, ['published @scope/pkg-linux-x64@1.0.0-rc.1 (rc)']);
});

test('publishPackage skips a package that is already published', () => {
    const { runner, logs, deps } = setup({
        view: { status: 0, stdout: '1.0.0-rc.1\n' },
    });

    publishPackage(PKG_DIR, 'rc', deps);

    assert.ok(!runner.calls.some((c) => c.args[0] === 'publish'), 'must not publish');
    assert.deepStrictEqual(logs, ['skipped @scope/pkg-linux-x64@1.0.0-rc.1 (already published)']);
});

test('publishPackage fails without publishing when the existence query errors', () => {
    const { runner, deps } = setup({
        view: { status: 1, stderr: 'npm error code ETIMEDOUT\nnpm error network timeout' },
        publish: { status: 0 },
    });

    assert.throws(() => publishPackage(PKG_DIR, 'rc', deps), /@scope\/pkg-linux-x64@1\.0\.0-rc\.1/);
    assert.ok(!runner.calls.some((c) => c.args[0] === 'publish'), 'must not publish');
});

test('publishPackage fails without publishing when npm cannot be started', () => {
    const { runner, deps } = setup({
        view: { status: null, error: new Error('spawn npm ENOENT') },
        publish: { status: 0 },
    });

    assert.throws(() => publishPackage(PKG_DIR, 'rc', deps));
    assert.ok(!runner.calls.some((c) => c.args[0] === 'publish'), 'must not publish');
});

test('publishPackage fails when npm publish fails', () => {
    const { deps } = setup({
        view: { status: 1, stderr: 'npm error code E404' },
        publish: { status: 1, stderr: 'npm error provenance generation failed' },
    });

    assert.throws(() => publishPackage(PKG_DIR, 'rc', deps), (err) => {
        assert.match(err.message, /@scope\/pkg-linux-x64@1\.0\.0-rc\.1/);
        assert.ok(!err.message.includes('\n'), 'message must be a single line');
        return true;
    });
});

// --- Release notes (US-6) ---

const ROOT = path.join('fake', 'repo');
const notesFile = (tag) => path.join(ROOT, '.github', 'release-notes', `${tag}.md`);

test('notesPath returns the notes file path when it has content', () => {
    const readFile = fakeReader({ [notesFile('v1.0.0')]: '# 1.0.0\n\nMigration notes.\n' });
    assert.strictEqual(notesPath('v1.0.0', { root: ROOT, readFile }), notesFile('v1.0.0'));
});

test('notesPath returns nothing when the notes file is missing', () => {
    const readFile = fakeReader({});
    assert.strictEqual(notesPath('v1.0.0', { root: ROOT, readFile }), null);
});

test('notesPath returns nothing when the notes file is whitespace only', () => {
    const readFile = fakeReader({ [notesFile('v1.0.0')]: '  \n\t\n' });
    assert.strictEqual(notesPath('v1.0.0', { root: ROOT, readFile }), null);
});

// --- Native package versions (release-version-optional-deps US-1, US-2) ---

const NATIVE = [
    '@jeancodogno/specforce-kit-linux-x64',
    '@jeancodogno/specforce-kit-linux-arm64',
    '@jeancodogno/specforce-kit-darwin-x64',
    '@jeancodogno/specforce-kit-darwin-arm64',
    '@jeancodogno/specforce-kit-win32-x64',
];

function nativeDeps(version, extra = {}) {
    const deps = {};
    for (const name of NATIVE) {
        deps[name] = version;
    }
    return { ...deps, ...extra };
}

function assertSingleLine(err) {
    assert.ok(!err.message.includes('\n'), 'message must be a single line');
}

test('NATIVE_PACKAGES lists the five Makefile platforms', () => {
    assert.deepStrictEqual(NATIVE_PACKAGES, NATIVE);
});

test('checkNativeVersions passes when all five native packages are at the version', () => {
    assert.doesNotThrow(() => checkNativeVersions(nativeDeps('1.0.0'), '1.0.0'));
});

test('checkNativeVersions fails naming a package at an old version and both versions', () => {
    const deps = nativeDeps('1.0.0', { '@jeancodogno/specforce-kit-darwin-x64': '0.9.0' });
    assert.throws(
        () => checkNativeVersions(deps, '1.0.0'),
        (err) => {
            assert.match(err.message, /specforce-kit-darwin-x64/);
            assert.match(err.message, /0\.9\.0/);
            assert.match(err.message, /1\.0\.0/);
            assert.ok(!/linux-x64|win32-x64/.test(err.message), 'only the wrong package is named');
            assertSingleLine(err);
            return true;
        }
    );
});

test('checkNativeVersions fails naming a missing package', () => {
    const deps = nativeDeps('1.0.0');
    delete deps['@jeancodogno/specforce-kit-win32-x64'];
    assert.throws(
        () => checkNativeVersions(deps, '1.0.0'),
        (err) => {
            assert.match(err.message, /missing/);
            assert.match(err.message, /specforce-kit-win32-x64/);
            assertSingleLine(err);
            return true;
        }
    );
});

test('checkNativeVersions fails naming an unexpected specforce-kit package', () => {
    const deps = nativeDeps('1.0.0', { '@jeancodogno/specforce-kit-foo': '1.0.0' });
    assert.throws(
        () => checkNativeVersions(deps, '1.0.0'),
        (err) => {
            assert.match(err.message, /unexpected/);
            assert.match(err.message, /specforce-kit-foo/);
            assertSingleLine(err);
            return true;
        }
    );
});

test('checkNativeVersions ignores unrelated optional dependencies', () => {
    const deps = nativeDeps('1.0.0', { fsevents: '^2.3.0' });
    assert.doesNotThrow(() => checkNativeVersions(deps, '1.0.0'));
});

test('checkNativeVersions without optionalDependencies names all five packages', () => {
    assert.throws(
        () => checkNativeVersions(undefined, '1.0.0'),
        (err) => {
            for (const name of NATIVE) {
                assert.ok(err.message.includes(name), `${name} must be named`);
            }
            assertSingleLine(err);
            return true;
        }
    );
});

test('checkRelease passes when versions and native packages agree', () => {
    assert.doesNotThrow(() =>
        checkRelease('v1.0.0', {
            packageVersion: '1.0.0',
            builtinVersion: '1.0.0',
            optionalDependencies: nativeDeps('1.0.0'),
        })
    );
});

test('checkRelease reports a native package mismatch as one line', () => {
    assert.throws(
        () =>
            checkRelease('v1.0.0', {
                packageVersion: '1.0.0',
                builtinVersion: '1.0.0',
                optionalDependencies: nativeDeps('1.0.0', {
                    '@jeancodogno/specforce-kit-linux-arm64': '0.9.0',
                }),
            }),
        (err) => {
            assert.match(err.message, /specforce-kit-linux-arm64/);
            assertSingleLine(err);
            return true;
        }
    );
});

test('checkRelease still reports a version mismatch first', () => {
    assert.throws(
        () =>
            checkRelease('v1.0.0', {
                packageVersion: '1.0.0-rc.1',
                builtinVersion: '1.0.0',
                optionalDependencies: nativeDeps('1.0.0'),
            }),
        /version mismatch/
    );
});

function pkgText(obj) {
    return JSON.stringify(obj, null, 2) + '\n';
}

test('syncNativeVersions sets all five native packages to the package version', () => {
    const input = pkgText({ name: 'x', version: '1.1.0', optionalDependencies: nativeDeps('1.0.0') });
    const out = JSON.parse(syncNativeVersions(input));
    assert.deepStrictEqual(out.optionalDependencies, nativeDeps('1.1.0'));
});

test('syncNativeVersions adds a missing native package', () => {
    const deps = nativeDeps('1.1.0');
    delete deps['@jeancodogno/specforce-kit-darwin-arm64'];
    const out = JSON.parse(syncNativeVersions(pkgText({ version: '1.1.0', optionalDependencies: deps })));
    assert.strictEqual(out.optionalDependencies['@jeancodogno/specforce-kit-darwin-arm64'], '1.1.0');
    assert.doesNotThrow(() => checkNativeVersions(out.optionalDependencies, '1.1.0'));
});

test('syncNativeVersions adds optionalDependencies when absent', () => {
    const out = JSON.parse(syncNativeVersions(pkgText({ version: '1.1.0' })));
    assert.deepStrictEqual(out.optionalDependencies, nativeDeps('1.1.0'));
});

test('syncNativeVersions removes an extra specforce-kit package and keeps unrelated entries', () => {
    const deps = nativeDeps('1.1.0', { '@jeancodogno/specforce-kit-foo': '1.0.0', fsevents: '^2.3.0' });
    const out = JSON.parse(syncNativeVersions(pkgText({ version: '1.1.0', optionalDependencies: deps })));
    assert.ok(!('@jeancodogno/specforce-kit-foo' in out.optionalDependencies));
    assert.strictEqual(out.optionalDependencies.fsevents, '^2.3.0');
});

test('syncNativeVersions keeps unrelated fields of package.json', () => {
    const input = { name: 'x', version: '1.1.0', bin: { x: 'index.js' }, license: 'MIT' };
    const out = JSON.parse(syncNativeVersions(pkgText(input)));
    assert.strictEqual(out.name, 'x');
    assert.deepStrictEqual(out.bin, { x: 'index.js' });
    assert.strictEqual(out.license, 'MIT');
});

test('syncNativeVersions is idempotent', () => {
    const input = pkgText({ version: '1.1.0', optionalDependencies: nativeDeps('1.0.0') });
    const once = syncNativeVersions(input);
    assert.strictEqual(syncNativeVersions(once), once);
});

test('syncNativeVersions returns already synced text unchanged', () => {
    const input = pkgText({ name: 'x', version: '1.1.0', optionalDependencies: nativeDeps('1.1.0') });
    assert.strictEqual(syncNativeVersions(input), input);
});

test('syncNativeVersions writes 2-space JSON with a trailing newline', () => {
    const out = syncNativeVersions(JSON.stringify({ version: '1.1.0' }));
    assert.ok(out.endsWith('}\n'));
    assert.ok(!out.endsWith('\n\n'));
    assert.match(out, /^\{\n  "version": "1\.1\.0",\n  "optionalDependencies": \{\n    "@jeancodogno/);
});
