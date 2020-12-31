db = db.getSiblingDB('statistics')

db.createUser(
	{
		user: "admin",
		pwd: "devpass",
		roles: [
			{
				role: "readWrite",
				db: "statistics"
			}
		]
	},
	{w: "majority", wtimeout: 5000}
);