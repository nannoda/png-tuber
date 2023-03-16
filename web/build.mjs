import * as esbuild from 'esbuild'

// get args
let watch = false
let serve = false
const args = process.argv.slice(2)
for (const arg of args) {
    if (arg == "watch") watch = true
    if (arg == "serve") serve = true
}

const buildCtx = await esbuild.context({
    entryPoints: ['index.ts'],
    bundle: true,
    sourcemap: true,
    outfile: 'index.js',
}
)

if (watch) {
    await buildCtx.watch()
}
if (serve) {
    await buildCtx.serve(
        {
            servedir: '.',
            port: 8080
        }
    )
}
