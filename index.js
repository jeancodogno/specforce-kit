#!/usr/bin/env node

const { spawnSync } = require('child_process');
const path = require('path');

/**
 * SpecForce Kit Proxy
 * Detects the current platform/arch and executes the corresponding native binary.
 */

const platform = process.platform;
const arch = process.arch;
const pkgName = `@jeancodogno/specforce-kit-${platform}-${arch}`;

let binaryPath;

try {
    // 1. Try to resolve the platform-specific package
    const pkgJsonPath = require.resolve(`${pkgName}/package.json`);
    const pkgDir = path.dirname(pkgJsonPath);
    const pkgJson = require(pkgJsonPath);
    
    // 2. Extract binary path from package.json bin field
    const binRelativePath = typeof pkgJson.bin === 'string' ? pkgJson.bin : pkgJson.bin.specforce;
    binaryPath = path.resolve(pkgDir, binRelativePath);
} catch (e) {
    // Fallback 1: Look for a binary in the root (built via 'make build' or 'prepare')
    const localBinaryName = platform === 'win32' ? 'specforce.exe' : 'specforce';
    const rootBinaryPath = path.join(__dirname, localBinaryName);
    const fs = require('fs');

    if (fs.existsSync(rootBinaryPath)) {
        binaryPath = rootBinaryPath;
    } else {
        // Fallback 2: Local development folder structure
        const localPkgDir = path.join(__dirname, 'npm', `${platform}-${arch}`);
        try {
            const pkgJson = require(path.join(localPkgDir, 'package.json'));
            const binRelativePath = typeof pkgJson.bin === 'string' ? pkgJson.bin : pkgJson.bin.specforce;
            binaryPath = path.resolve(localPkgDir, binRelativePath);
        } catch (err) {
            runDiagnostic(platform, arch, pkgName);
            process.exit(1);
        }
    }
}

/**
 * Runs a diagnostic check to help the user fix their environment.
 */
function runDiagnostic(platform, arch, pkgName) {
    const fs = require('fs');
    
    // Aesthetic Colors
    const RESET = "\x1b[0m";
    const BOLD = "\x1b[1m";
    const RED = "\x1b[38;5;196m";
    const MINT = "\x1b[38;5;48m";
    const GRAY = "\x1b[38;5;245m";

    console.error(`\n${BOLD}${RED}┌── [ SpecForce Diagnostic ] ──────────────────────────┐${RESET}`);
    console.error(`${RED}│${RESET}  ${BOLD}Error:${RESET} Native binary for ${platform}-${arch} not found.     ${RED}│${RESET}`);
    console.error(`${RED}└──────────────────────────────────────────────────────┘${RESET}`);

    try {
        const npmPrefix = spawnSync('npm', ['config', 'get', 'prefix'], { encoding: 'utf8' }).stdout.trim();
        const binDir = platform === 'win32' ? npmPrefix : path.join(npmPrefix, 'bin');
        
        console.error(`\n${BOLD}Target Package:${RESET} ${pkgName}`);
        console.error(`${BOLD}Expected Bin Directory:${RESET} ${binDir}`);

        if (!process.env.PATH.includes(binDir)) {
            console.error(`\n${BOLD}${MINT}>>> CAUSE DETECTED: PATH MISCONFIGURATION <<<${RESET}`);
            console.error(`${GRAY}Your npm global bin directory is not in your system's PATH.${RESET}\n`);
            
            if (platform === 'win32') {
                console.error(`To fix this on Windows, run this in PowerShell as Administrator:`);
                console.error(`  ${BOLD}setx /M PATH "%PATH%;${binDir}"${RESET}`);
            } else {
                const shellProfile = process.env.SHELL?.includes('zsh') ? '~/.zshrc' : '~/.bashrc';
                console.error(`To fix this on macOS/Linux, add this to your ${shellProfile}:`);
                console.error(`  ${BOLD}export PATH="${binDir}:$PATH"${RESET}`);
                console.error(`\nThen reload your shell: ${BOLD}source ${shellProfile}${RESET}`);
            }
        } else {
            console.error(`\n${BOLD}${MINT}>>> CAUSE DETECTED: MISSING BINARY <<<${RESET}`);
            console.error(`${GRAY}The binary was not found even though your PATH seems correct.${RESET}\n`);
            console.error(`If you are a developer, run: ${BOLD}make build${RESET}`);
            console.error(`Otherwise, try: ${BOLD}npm install -g @jeancodogno/specforce-kit --force${RESET}`);
        }
    } catch (e) {
        console.error(`\n${GRAY}Failed to run full diagnostic. Please ensure 'npm' is installed.${RESET}`);
    }
    console.error("");
}

// 3. Execute the binary forwarding all arguments and inheriting stdio
const result = spawnSync(binaryPath, process.argv.slice(2), {
    stdio: 'inherit'
});

// 4. Propagate the exit code
process.exit(result.status ?? 0);
