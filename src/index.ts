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
  return new Date(`${day}, ${time} GMT`);
}

function toLocalizedTimestamp(date: Date): string {
  const options = {
    timeZone: 'America/Chicago',
    timeZoneName: 'short',
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  };
  return date.toLocaleString('en-US', options);
}

interface IRawData {
  [key: string]: boolean,
};

interface IReading {
  date: Date,
  localizedTimestamp: string,
  isDry: boolean,
};

function getCurrentReading(data: IRawData): void | IReading {
  const sortedKeys = Object.keys(data).sort((a, b) => a.localeCompare(b));
  if (!sortedKeys.length) return undefined;

  const key = sortedKeys[sortedKeys.length - 1];
  const date = getDate(key);
  const localizedTimestamp = toLocalizedTimestamp(date);
  const isDry = data[key];
  return {
    date,
    localizedTimestamp,
    isDry,
  };
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
      const currentReading = getCurrentReading(data);
      console.log({ currentReading });
    },
    err => {
      console.log('error: ', err.message);
      process.exit();
    },
  );
}

main();




