'use client'

import { useState } from 'react'

const shells = [
  {
    name: 'Bash',
    code: `# Add to ~/.bashrc or ~/.bash_profile
eval "$(kasoku shell-hook bash)"

# That's it! Now all your commands are automatically cached
npm run build
go build ./...
cargo test`
  },
  {
    name: 'Zsh',
    code: `# Add to ~/.zshrc
eval "$(kasoku shell-hook zsh)"

# Kasoku now transparently caches your builds
npm run build
go build ./...
cargo test`
  },
  {
    name: 'Fish',
    code: `# Add to ~/.config/fish/config.fish
kasoku shell-hook fish | source

# Every build command is now intelligently cached
npm run build
go build ./...
cargo test`
  }
]

export function ShellTabs() {
  const [activeTab, setActiveTab] = useState(0)

  return (
    <div className="max-w-3xl mx-auto">
      {/* Tab Headers */}
      <div className="flex border-b border-gray-700">
        {shells.map((shell, index) => (
          <button
            key={shell.name}
            onClick={() => setActiveTab(index)}
            className={`px-6 py-3 font-mono text-sm transition-colors ${
              activeTab === index
                ? 'bg-gray-800 text-white border-b-2 border-green-400'
                : 'text-gray-400 hover:text-white hover:bg-gray-800/50'
            }`}
          >
            {shell.name}
          </button>
        ))}
      </div>

      {/* Tab Content */}
      <div className="bg-gray-900 text-white rounded-b-lg p-6 font-mono text-sm">
        <pre className="whitespace-pre-wrap">
          <code className="text-gray-400">{shells[activeTab].code.split('\n')[0]}</code>
          {'\n'}
          <code className="text-green-400">{shells[activeTab].code.split('\n')[1]}</code>
          {'\n\n'}
          <code className="text-gray-400">{shells[activeTab].code.split('\n')[3]}</code>
          {'\n'}
          <code className="text-blue-400">{shells[activeTab].code.split('\n')[4]}</code>
          {'\n'}
          <code className="text-blue-400">{shells[activeTab].code.split('\n')[5]}</code>
          {'\n'}
          <code className="text-blue-400">{shells[activeTab].code.split('\n')[6]}</code>
        </pre>
      </div>

      {/* Benefits */}
      <div className="mt-6 grid md:grid-cols-3 gap-4 text-sm">
        <div className="bg-gray-50 rounded-lg p-4">
          <div className="font-semibold mb-1">✨ Zero Friction</div>
          <div className="text-muted-foreground">No need to wrap commands with kasoku exec</div>
        </div>
        <div className="bg-gray-50 rounded-lg p-4">
          <div className="font-semibold mb-1">🎯 Smart Detection</div>
          <div className="text-muted-foreground">Automatically recognizes build commands</div>
        </div>
        <div className="bg-gray-50 rounded-lg p-4">
          <div className="font-semibold mb-1">⚡ Instant Speed</div>
          <div className="text-muted-foreground">Cache kicks in transparently</div>
        </div>
      </div>
    </div>
  )
}
