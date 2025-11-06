const { execSync } = require('child_process');
const os = require('os');

const isWindows = os.platform() === 'win32';

if (isWindows) {
  execSync('xcopy /E /I ..\\storybook-shared\\dist .\\storybook-shared-dist', {
    stdio: 'inherit',
  });
} else {
  execSync('cp -R ../storybook-shared/dist ./storybook-shared-dist', {
    stdio: 'inherit',
  });
}
