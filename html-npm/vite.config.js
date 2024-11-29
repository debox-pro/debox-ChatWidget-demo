import { defineConfig } from 'vite'

export default defineConfig({
    base: './',  // Set to relative path
    root: 'src',  // Set Vite's root directory to the src folder
    build: {
        outDir: '../dist',  // Set the build output directory to the dist folder under the root directory
        emptyOutDir: true,  // Force emptying the output directory
    }
})
