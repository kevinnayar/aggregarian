import * as admin from 'firebase-admin';
import { database, ServiceAccount } from 'firebase-admin';
import * as serviceAccount from './service-account-key.json';

function getDatabaseInstance(credentials: ServiceAccount, projectId: string): database.Database {
  admin.initializeApp({
    credential: admin.credential.cert(credentials),
    databaseURL: `https://${projectId}.firebaseio.com`,
  });
  return admin.database();
}

function getUTCDateFromKey(key: string): Date {
  const [date, time] = key.split('_');
  return new Date(`${date}, ${time} GMT`);
}

function localeStringToFormatted(localeString: string): string {
  const [date, time] = localeString.split(',');
  const [month, day, year] = date.split('/');
  return `${year}-${month}-${day},${time}`;
}

function dateToLocalDateString(date: Date, locale: string, tz: string): string {
  const options = {
    timeZone: tz,
    timeZoneName: 'short',
    hour12: false,
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  };
  const localeString = date.toLocaleString(locale, options);
  const formatted = localeStringToFormatted(localeString);
  return formatted;
}

interface IRawData {
  [key: string]: boolean,
};

interface IReading {
  dateUTC: Date,
  localDateString: string,
  locale: string,
  timezone: string,
  isDry: boolean,
};

function getLatestReading(data: IRawData, locale: string, tz: string): void | IReading {
  const sortedKeys = Object.keys(data).sort((a, b) => a.localeCompare(b));
  if (!sortedKeys.length) return undefined;

  const key = sortedKeys[sortedKeys.length - 1];
  const dateUTC = getUTCDateFromKey(key);
  const localDateString = dateToLocalDateString(dateUTC, locale, tz);
  const isDry = data[key];
  return {
    dateUTC,
    localDateString,
    locale,
    timezone: tz,
    isDry,
  };
}

function main() {
  const LOCALE = 'en-US';
  const TIMEZONE = 'America/Chicago';
  const projectId = serviceAccount.project_id;
  const credentials = serviceAccount as ServiceAccount;
  
  const db = getDatabaseInstance(credentials, projectId);
  const ref = db.ref();

  ref.on(
    'value',
    snapshot => {
      const data = snapshot.val();
      const latestReading = getLatestReading(data, LOCALE, TIMEZONE);
      console.log({ latestReading });
    },
    err => {
      console.log('error: ', err.message);
      process.exit();
    },
  );
}

main();




