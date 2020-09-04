import * as admin from 'firebase-admin';
import { database, ServiceAccount } from 'firebase-admin';
import * as serviceAccount from '../config/service-account-key.json';

function getDatabaseInstance(credentials: ServiceAccount, projectId: string): database.Database {
  admin.initializeApp({
    credential: admin.credential.cert(credentials),
    databaseURL: `https://${projectId}.firebaseio.com`,
  });
  return admin.database();
}

interface IRawDataLogEntry {
  ReadableDate: string,
  IsDry: boolean,
};

interface IRawDataLog {
  [key: string]: IRawDataLogEntry,
};

interface IRawData {
  log: IRawDataLog,
};

interface IReading {
  timestampUtcIso8601: string,
  timestampLocalReadable: string,
  isDry: boolean,
}

function formatLogEntry(key: string, log: IRawDataLog): IReading {
  const entry: IRawDataLogEntry = log[key];
  if (!entry) throw new Error(`Could not find item: ${key} in log`);

  return {
    timestampUtcIso8601: key,
    timestampLocalReadable: entry.ReadableDate,
    isDry: entry.IsDry,
  };
}

function getLogFromData(data: IRawData): IRawDataLog {
  if (!data.log) throw new Error('Could not find log in data');
  return data.log;
}

function getSortedLogKeys(log: IRawDataLog): string[] {
  const keys = Object.keys(log).sort((a, b) => b.localeCompare(a));
  if (!keys.length) throw new Error('Could not find any keys in log');
  return keys;
}

function getNLatestReadings(data: IRawData, n: number): IReading[] {
  const log = getLogFromData(data);
  const keys = getSortedLogKeys(log);
  const nKeys = keys.slice(0, n);

  const readings = nKeys.reduce((all: IReading[], key: string) => {
    const reading = formatLogEntry(key, log);
    all.push(reading);
    return all;
  }, []);
  return readings;
}

function getLatestReading(data: IRawData): void | IReading {
  const log = getLogFromData(data);
  const keys = getSortedLogKeys(log);

  const reading = formatLogEntry(keys[0], log);
  return reading;
}

function main() {
  const projectId = serviceAccount.project_id;
  const credentials = serviceAccount as ServiceAccount;
  
  const db = getDatabaseInstance(credentials, projectId);
  const ref = db.ref();

  ref.on(
    'value',
    snapshot => {
      const data = snapshot.val();

      const lastReading = getLatestReading(data);
      console.log('lastReading', lastReading);

      const last10Readings = getNLatestReadings(data, 10);
      console.log('last10Readings', last10Readings);
      process.exit();
    },
    err => {
      console.log('error: ', err.message);
      process.exit();
    },
  );
}

main();




