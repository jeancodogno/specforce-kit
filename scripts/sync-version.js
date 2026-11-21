const fs = require('fs');
const path = require('path');
const { syncNativeVersions } = require('./release-tools');

const packageJsonPath = path.join(__dirname, '..', 'package.json');
const constantsGoPath = path.join(__dirname, '..', 'src', 'internal', 'core', 'constants.go');

try {
    const packageJsonText = fs.readFileSync(packageJsonPath, 'utf8');
    const packageJson = JSON.parse(packageJsonText);
    const version = packageJson.version;

    if (!version) {
        console.error('Error: No version found in package.json');
        process.exit(1);
    }

    const syncedPackageJson = syncNativeVersions(packageJsonText);
    if (syncedPackageJson !== packageJsonText) {
        fs.writeFileSync(packageJsonPath, syncedPackageJson);
        console.log(`Successfully updated native package versions to "${version}" in package.json`);
    }

    let constantsGo = fs.readFileSync(constantsGoPath, 'utf8');
    
    // Support both const and var declarations during the transition
    const versionRegex = /(const|var)\s+Version\s+=\s+"[^"]+"/;
    const newDeclaration = `var Version = "${version}"`;

    if (versionRegex.test(constantsGo)) {
        constantsGo = constantsGo.replace(versionRegex, newDeclaration);
        fs.writeFileSync(constantsGoPath, constantsGo);
        console.log(`Successfully updated Version to "${version}" in constants.go`);
    } else {
        console.error('Error: Could not find Version declaration in constants.go');
        process.exit(1);
    }
} catch (error) {
    console.error(`Error: ${error.message}`);
    process.exit(1);
}
