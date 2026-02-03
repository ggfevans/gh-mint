# gh-repo-defaults demo
# Run: vhs static/demo.tape

Output static/hero.gif

Set Shell bash
Set FontSize 16
Set Width 900
Set Height 550
Set Padding 20
Set Theme "Catppuccin Mocha"
Set TypingSpeed 40ms
Set CursorBlink false

# Show what this thing does
Type "gh-repo-defaults --help"
Enter
Sleep 2s

# List available profiles
Type "gh-repo-defaults profiles list"
Enter
Sleep 2.5s

# Inspect a profile
Type "gh-repo-defaults profiles show oss"
Enter
Sleep 3s

# Show the create command
Type "gh-repo-defaults create --help"
Enter
Sleep 2s

# One-liner example
Type "# One command. Fully configured repo."
Enter
Sleep 500ms
Type "gh-repo-defaults create my-project --profile oss --public"
Sleep 2s
