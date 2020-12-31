db = db.getSiblingDB('tweets')

db.createUser(
	{
		user: "admin",
		pwd: "devpass",
		roles: [
			{
				role: "readWrite",
				db: "tweets"
			}
		]
	},
	{w: "majority", wtimeout: 5000}
);