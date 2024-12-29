db.createUser({
    user: "development",
    pwd: "testpassword",
    roles: [
      {
        role: "readWrite",
        db: "admin"
      }
    ]
  });