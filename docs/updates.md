# Updates Summary

## Latest Fixes (Jan 6, 2026 - Final)

### 1. Clean TUI (No Log Text in UI)
- **Issue**: Log text like "SMTP server listening on port 2525" and "SMTP server stopped" appeared in TUI when pressing SPACE or receiving emails
- **Fix**: Removed all `log.Printf()` statements from SMTP server and email handling
- **Result**: TUI is now completely clean with no unwanted text appearing
- **Note**: Only `log.Fatal()` kept for critical startup errors (before TUI starts)

### 2. Correct Initial Server State
- **Issue**: TUI showed "Stopped" status on startup even though server was actually running
- **Fix**: Start server synchronously before creating GUI and drawing initial layout
- **Result**: Initial state now correctly shows "Running" (in green) when app starts

### 3. SPACE Key Toggle Fixed

### 1. SPACE Key Toggle Fixed
- **Issue**: SPACE key wasn't toggling the server
- **Fix**: Changed from character literal `' '` to `gocui.KeySpace` constant
- **Result**: SPACE now properly starts/stops SMTP server with immediate feedback

### 2. Real-Time Email Updates Fixed
- **Issue**: Emails weren't appearing in real-time when received
- **Fix**: Wrapped UI updates in `g.Update()` callback in the notification goroutine
- **Result**: New emails appear instantly when received from Laravel or any other mailer

### 3. Email List Display Changed to Receiver
- **Before**: Showed `sender | subject` (From address)
- **After**: Shows `receiver | subject` (To address)
- **Reason**: More useful for testing - see who email was sent to
- **Format**: `> user@example.com | Test Subject...`

### 4. Removed Background from First Email
- **Issue**: First email always had blue background
- **Fix**: Removed `v.Highlight = true` and highlight color settings from email view
- **Result**: Only selected email shows blue background with `>` prefix

### 5. GUI Update Calls Fixed
- Changed all keybinding handlers to use `gui` parameter name
- Ensures `g.Update()` is called correctly for all UI changes
- Fixed SPACE key callback to properly refresh server info panel

## Completed Changes

### 1. TUI Layout Improvements
- **Wider email list pane**: Increased from 30 to 45 characters
- **No line wrapping**: Disabled wrap on email list and server info
- **Text truncation**: Email subjects and From addresses are truncated at 35 chars with "..." suffix
- **Better proportions**: Server panel takes 1/3 height, email list takes remaining space

### 2. Visual Enhancements
- **Colors added**:
  - Cyan: Frame and titles for left panels
  - Green: Frame and title for main panel
  - Blue: Selected email background with white text
  - Red: "Stopped" status
  - Green: "Running" status
  - Yellow: Action prompts
- **Highlight**: Selected emails show blue background with white text
- **Colored headers**: Email detail view has cyan labels for ID, From, To, Subject, Date, Body

### 3. Keyboard Controls
- **j/k**: Navigate down/up through emails
- **d**: Delete selected email
- **SPACE**: Toggle SMTP server on/off
- **ESC**: Go back to homepage (deselects email)
- **q**: Quit application
- **Ctrl+C**: Quit application (new)

### 4. Auto-Refresh on New Emails
- **Channel-based notification**: When new email arrives, TUI automatically refreshes
- **Real-time updates**: Email list and server info update immediately when email received
- **No manual refresh needed**: j/k navigation will always show latest emails

### 5. Server State Management
- **Reliable tracking**: Added `running` flag to SMTPServer struct
- **Improved Toggle**: SPACE now reliably starts/stops server with immediate UI feedback
- **State persistence**: Server state tracked in memory, not via network connection

### 6. Build Changes
- **Build output**: All binaries now go to `build/` directory
- **Platform binaries**:
  - `build/lazysmtp` (current platform)
  - `build/lazysmtp-linux-amd64`
  - `build/lazysmtp-darwin-amd64`
  - `build/lazysmtp-darwin-arm64`
  - `build/lazysmtp-windows-amd64.exe`

## File Changes

### src/types.go
- Added `NewEmailChan chan struct{}` to AppState for real-time notifications

### src/smtp.go
- Removed all `log.Printf()` statements (prevents text appearing in TUI)
- Removed `log.Println()` from `Stop()` method
- Removed unused `log` import
- Added `notify chan struct{}` to Backend
- Created `NewBackend()` constructor
- Added `notify chan struct{}` to Session
- Added notification send in `Data()` when email is saved
- Added `running bool` to SMTPServer
- Modified `Start()` to set `running = true`
- Modified `Stop()` to set `running = false`
- Modified `IsRunning()` to return `running` instead of checking network
- Removed unused `net` import

### src/tui.go
- Increased left panel width to 45
- Set server panel to 1/3 of screen height
- Disabled wrapping on all left panels
- Added frame colors (Cyan, Green)
- Removed `Highlight` from email view (fixes first email background issue)
- Added `Ctrl+C` keybinding
- Added `ESC` keybinding to return to homepage
- Changed SPACE keybinding to use `gocui.KeySpace` constant
- Improved keybinding for j/k to work even with 0 emails
- Fixed all keybinding handlers to use `gui` parameter for proper UI updates

### src/main.go
- Start server synchronously before GUI creation (ensures correct initial state)
- Removed `log.Printf()` statements from email update goroutine
- Created notification channel for real-time email updates
- Updated AppState initialization with notification channel
- Added goroutine to listen for new emails and refresh UI using `g.Update()`
- Changed email list display from `From` to `To` (receiver)
- Updated email list truncation with "..."
- Added color codes for status, labels, and controls
- Updated controls help text

## Testing with Laravel

To test with the Laravel app at `/home/mouayed/Desktop/dev/php/laravel-playground`:

1. Start lazySMTP:
   ```bash
   ./build/lazysmtp
   ```

2. In another terminal, send test email:
   ```bash
   cd /home/mouayed/Desktop/dev/php/laravel-playground
   php artisan mailtest user@example.com
   ```

The email should appear immediately in lazySMTP without needing to press any keys.

## Known Limitations

- GUI requires a TTY (cannot run in background or via systemd directly)
- For headless/server environments, consider adding a CLI mode (not yet implemented)
