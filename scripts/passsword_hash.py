import bcrypt
passwords = {
    "ADMIN_PASS": "AdminPassword123!",
    "CLIENT_PASS": "ClientPassword123!",
}
for label, pwd in passwords.items():
    hashed = bcrypt.hashpw(pwd.encode(), bcrypt.gensalt()).decode()
    print(f"{label}={hashed}")