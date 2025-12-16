import { existsSync, mkdirSync, readFileSync } from "node:fs"
import { execSync } from "node:child_process"

const outputDirs = {
  windows: 'examples/.terraform/providers/registry.terraform.io/ashleyjackson/instatus/1.0.0/windows_amd64/',
  linux: 'examples/.terraform/providers/registry.terraform.io/ashleyjackson/instatus/1.0.0/linux_amd64/',
  macos: 'examples/.terraform/providers/registry.terraform.io/ashleyjackson/instatus/1.0.0/darwin_amd64/',
}

const args = process.argv.slice(2)
if (args.length !== 1 || !['windows', 'linux', 'macos'].includes(args[0])) {
  console.error('Usage: ts-node build.ts <windows|linux|macos>')
  process.exit(1)
}

const platform = args[0] as keyof typeof outputDirs
const outputDir = outputDirs[platform]

// Ensure output directory exists
if (!existsSync(outputDir)) {
  mkdirSync(outputDir, { recursive: true })
}

if (platform === 'windows') {
  console.log('Building for Windows...')
  console.log(`Copying to ${outputDir}...`)
  execSync(`set GOOS=windows&& set GOARCH=amd64&& go build -o "${outputDir}terraform-provider-instatus.exe"`, { stdio: 'inherit' })
}
if (platform === 'linux') {
  console.log('Building for Linux... (Untested)')
  console.log(`Copying to ${outputDir}...`)
  execSync(`GOOS=linux GOARCH=amd64 go build -o "${outputDir}terraform-provider-instatus"`, { stdio: 'inherit' })
}
if (platform === 'macos') {
  console.log('Building for macOS... (Untested)')
  console.log(`Copying to ${outputDir}...`)
  execSync(`GOOS=darwin GOARCH=amd64 go build -o "${outputDir}terraform-provider-instatus.dmg"`, { stdio: 'inherit' })
}
console.log('Build complete.')