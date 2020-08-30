import * as admin from 'firebase-admin';
import { database, ServiceAccount } from 'firebase-admin';
import * as serviceAccount from './service-account-key.json';

function getDatabase(credentials: ServiceAccount, projectId: string): database.Database {
  admin.initializeApp({
    credential: admin.credential.cert(credentials),
    databaseURL: `https://${projectId}.firebaseio.com`,
  });

  return admin.database();
}

function getDate(text: string): Date {
  const split = text.split('_');
  const day = split[0];
  const time = split[1];
  return new Date(`${day}, ${time}`);
}

interface IRawData {
  [key: string]: boolean,
};

interface IReading {
  date: Date,
  isDry: boolean,
};

function getReadings(data: IRawData): IReading[] {
  const readings: IReading[] = [];
  const sortedKeys = Object.keys(data).sort((a, b) => a.localeCompare(b));
  
  for (const key of sortedKeys) {
    const date = getDate(key);
    const isDry = data[key];
    readings.push({
      date,
      isDry,
    });
  }

  return readings;
}

function main() {
  const projectId = serviceAccount.project_id;
  const credentials = serviceAccount as ServiceAccount;
  
  const db = getDatabase(credentials, projectId);
  const ref = db.ref();

  ref.on(
    'value',
    snapshot => {
      const data = snapshot.val();
      const readings = getReadings(data);
      console.log({ readings });
    },
    err => {
      console.log('error: ', err.message);
      process.exit();
    },
  );
}

main();




