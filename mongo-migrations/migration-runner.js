import { MongoClient } from "mongodb";
import migration0001 from "./migrations/0001-add-unique-index-to-users.js";
import migration0002 from "./migrations/0002-add-unique-index-to-companies.js";
import dotenv from "dotenv";

dotenv.config();

const uri = process.env.MONGODB_URI;
const dbName = "mx-auth";
const MIGRATIONS_COLLECTION = "migrations";

// TODO walk migrations dir
const migrations = [
  { id: "0001-add-unique-index-to-users", func: migration0001 },
  { id: "0002-add-unique-index-to-companies", func: migration0002 },
];

async function runMigrations() {
  const client = new MongoClient(uri);

  try {
    await client.connect();
    const db = client.db(dbName);
    const migrationsCollection = db.collection(MIGRATIONS_COLLECTION);

    for (const migration of migrations) {
      const alreadyApplied = await migrationsCollection.findOne({
        id: migration.id,
      });

      if (!alreadyApplied) {
        await migration.func(db);
        await migrationsCollection.insertOne({
          id: migration.id,
          appliedAt: new Date(),
        });
        console.log(`Migration ${migration.id} applied.`);
      } else {
        console.log(`Migration ${migration.id} already applied. Skipping.`);
      }
    }
  } catch (err) {
    console.error("Migration failed:", err);
  } finally {
    await client.close();
  }
}

runMigrations();
