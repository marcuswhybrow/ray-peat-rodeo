import browserSync from 'browser-sync'
import chokidar from 'chokidar'
import build from './build.js'

chokidar
  .watch(['src', 'layouts'], {
    ignoreInitial: true
  })
  .on('ready', () => browserSync.init({
    host: 'localhost',
    port: 8000,
    server: './build',
    injectChanges: false,
    interval: 2000
  }))
  .on('all', async (...args) => {
    await build()
    browserSync.reload()
  });

(async function() {
  await build()
}())