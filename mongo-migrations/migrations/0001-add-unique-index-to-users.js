export default async function (db) {
  console.log("Running migration 0001: Creating unique index on users.username");
  const users = db.collection('users');
  await users.createIndex({ username: 1 }, { unique: true });
}