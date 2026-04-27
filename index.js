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
    // Fallback for local development if not in node_modules
    const localPkgDir = path.join(__dirname, 'npm', `${platform}-${arch}`);
    try {
        const pkgJson = require(path.join(localPkgDir, 'package.json'));
        const binRelativePath = typeof pkgJson.bin === 'string' ? pkgJson.bin : pkgJson.bin.specforce;
        binaryPath = path.resolve(localPkgDir, binRelativePath);
    } catch (err) {
        console.error(`[SpecForce] Error: Native binary for ${platform}-${arch} not found.`);
        console.error(`Requirement: The package "${pkgName}" must be installed.`);
        process.exit(1);
    }
}

// 3. Execute the binary forwarding all arguments and inheriting stdio
const result = spawnSync(binaryPath, process.argv.slice(2), {
    stdio: 'inherit'
});

// 4. Propagate the exit code
process.exit(result.status ?? 0);
