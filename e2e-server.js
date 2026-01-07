#!/usr/bin/env node

/**
 * Simple E2E test server
 * Serves the fullstack example for E2E testing
 */

const { spawn } = require('child_process');
const path = require('path');

const exampleDir = path.join(__dirname, 'examples', 'fullstack');

console.log('Starting E2E test server...');
console.log(`Example directory: ${exampleDir}`);

// Build frontend
console.log('\nBuilding frontend...');
const buildProcess = spawn('npm', ['run', 'build'], {
  cwd: exampleDir,
  stdio: 'inherit',
  shell: true
});

buildProcess.on('exit', (code) => {
  if (code !== 0) {
    console.error('Build failed');
    process.exit(1);
  }
  
  console.log('\nStarting server...');
  
  // Start Go server
  const serverProcess = spawn('go', ['run', 'main.go'], {
    cwd: exampleDir,
    stdio: 'inherit',
    shell: true,
    env: {
      ...process.env,
      PORT: '3000',
    }
  });
  
  // Handle shutdown
  process.on('SIGINT', () => {
    console.log('\nShutting down server...');
    serverProcess.kill();
    process.exit(0);
  });
  
  process.on('SIGTERM', () => {
    serverProcess.kill();
    process.exit(0);
  });
  
  serverProcess.on('exit', (code) => {
    process.exit(code || 0);
  });
});
