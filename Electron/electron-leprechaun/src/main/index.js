import { app, shell, BrowserWindow, ipcMain } from 'electron'
import { join } from 'path'
import { electronApp, optimizer, is } from '@electron-toolkit/utils'
import icon from '../../resources/icon.png?asset'
import { spawn } from 'child_process' // Node.js module for spawning processes
import waitOn from 'wait-on' // Make sure to install this with: npm install wait-on

let goServer

function createWindow() {
  // Create the browser window.
  const mainWindow = new BrowserWindow({
    width: 900,
    height: 670,
    show: false, // Don't show until Go server is ready
    autoHideMenuBar: true,
    ...(process.platform === 'linux' ? { icon } : {}),
    webPreferences: {
      preload: join(__dirname, '../preload/index.js'),
      sandbox: false,
      contextIsolation: false,
      webSecurity: false
    }
  })

  mainWindow.on('ready-to-show', () => {
    console.log('Window is ready to show')
    mainWindow.show() // Show the window once the server is ready
  })

  mainWindow.webContents.setWindowOpenHandler((details) => {
    shell.openExternal(details.url)
    return { action: 'deny' }
  })

  // Load the remote URL for development or the local html file for production.
  if (is.dev && process.env['ELECTRON_RENDERER_URL']) {
    console.log('Loading URL: ', process.env['ELECTRON_RENDERER_URL'])
    mainWindow.loadURL(process.env['ELECTRON_RENDERER_URL'])
  } else {
    console.log('Loading local index.html')
    mainWindow.loadFile(join(__dirname, '../renderer/index.html'))
  }
}

function startGoServer() {
  // Start the Go server as a child process
  goServer = spawn('go', ['run', '.'], { cwd: '../../backend' })

  // Listen for Go server output (stdout and stderr)
  goServer.stdout.on('data', (data) => {
    console.log(`Go server stdout: ${data.toString()}`) // Print the standard output
  })

  goServer.stderr.on('data', (data) => {
    console.error(`Go server stderr: ${data.toString()}`) // Print the error output
  })

  goServer.on('exit', (code, signal) => {
    console.log(`Go server exited with code ${code} and signal ${signal}`)
  })

  return goServer
}

app.whenReady().then(async () => {
  try {
    // Start the Go server
    goServer = startGoServer()

    // Wait for the Go server to be ready
    //await waitForGoServer()
    createWindow()

    // Once the server is ready, create the Electron window
  } catch (error) {
    console.error('Error starting Go server or Electron app:', error)
  }

  // Set app user model id for windows
  electronApp.setAppUserModelId('com.electron')

  // Default open or close DevTools by F12 in development
  app.on('browser-window-created', (_, window) => {
    optimizer.watchWindowShortcuts(window)
  })

  // IPC test
  ipcMain.on('ping', () => console.log('pong'))

  app.on('activate', function () {
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })
})

app.on('before-quit', () => {
  // Kill Go server process before quitting the app
  if (goServer) {
    console.log('Killing Go server...')
    goServer.kill('SIGKILL') // Use SIGKILL to forcefully terminate
  }

  // Ensure the Go process is killed by checking for lingering processes on port 8080
  const { exec } = require('child_process')

  console.log('Attempting to run lsof command...')
  exec('sh -c "lsof -t -i :8080"', (err, stdout, stderr) => {
    console.log('inside exec') // This will show if the exec is being called

    // Log error if there's an issue executing the command
    if (err) {
      console.error('Error executing lsof command:', err) // This will print more details if there's an error with exec
    }

    // Log stderr if any error occurred with lsof itself
    if (stderr) {
      console.error('stderr from lsof:', stderr)
    }

    // Log stdout to see the result of lsof
    if (stdout) {
      console.log('stdout from lsof:', stdout) // This will print the process IDs or an empty string
      console.log('Go server might still be running, killing processes on port 8080...')
      // Split the PIDs into an array and kill each one
      const pids = stdout.split('\n').filter(Boolean) // Split stdout into PIDs, and remove any empty entries
      pids.forEach((pid) => {
        exec(`kill -9 ${pid}`, (killErr, killStdout, killStderr) => {
          if (killErr) {
            console.error(`Error killing Go server process with PID ${pid}:`, killErr)
          } else {
            console.log(`Killed process with PID ${pid}`)
          }

          if (killStderr) {
            console.error('stderr from kill command:', killStderr) // Print stderr from kill command
          }
        })
      })
    } else {
      console.log('No processes found on port 8080')
    }
  })
})

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()
  }
})
