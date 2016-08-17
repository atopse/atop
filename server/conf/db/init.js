db.createUser({
    user:"server",
    pwd:"server", 
    roles: [
        { role:"readWrite", db:"atopServer" },
    ]
});