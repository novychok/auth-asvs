### Auth service powered by ASVS

#### V2: Authentication

**V2.1: Password Protection**  
✅ **2.1.1**: Password must be ≥12 chars after merging multiple spaces  
✅ **2.1.2**: Allow up to 128 chars, recommended ≤64  
✅ **2.1.3**: No password truncation; compress consecutive spaces
✅ **2.1.4**: Allow all printable Unicode (emojis, symbols)  
❌ **2.1.5**: Endpoint for password change is implemented  
❌ **2.1.6**: Current and new password required for password change
❌ **2.1.7**: Passwords checked against breach lists  
🎨 **2.1.8**: Frontend shows password strength meter  
🎨 **2.1.9**: No forced rules on char types (symbols, cases, etc.)  
✅ **2.1.10**: No periodic password change or password history requirements
🎨 **2.1.11**: Password pasting allowed in fields, and browser extensions/password managers are supported
🎨 **2.1.12**: User can temporarily reveal full or last entered password character

**V2.1: Basic Requirements**  
