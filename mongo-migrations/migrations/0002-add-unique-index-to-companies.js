export default async function (db) {
  console.log(
    "Running migration 0002: Creating unique index on companies.name"
  );
  const companies = db.collection("companies");
  await companies.createIndex({ name: 1 }, { unique: true });
}
