// const DB_NAME = "quicksaveCache";
// const STORE_NAME = "mainCache";
// const DB_VERSION = 1;

// // Open IndexedDB
// function openDB() {
//   return new Promise((resolve, reject) => {
//     const request = indexedDB.open(DB_NAME, DB_VERSION);

//     request.onupgradeneeded = (event) => {
//       const db = event.target.result;
//       if (!db.objectStoreNames.contains(STORE_NAME)) {
//         db.createObjectStore(STORE_NAME);
//       }
//     };

//     request.onsuccess = () => resolve(request.result);
//     request.onerror = () => reject(request.error);
//   });
// }

// // Function to write JSON data to IndexedDB
// export async function writeJSON(key, jsonData) {
//   const db = await openDB();
//   return new Promise((resolve, reject) => {
//     const transaction = db.transaction(STORE_NAME, "readwrite");
//     const store = transaction.objectStore(STORE_NAME);
//     const request = store.put(jsonData, key);

//     request.onsuccess = () => resolve(true);
//     request.onerror = () => reject(request.error);
//   });
// }

// // Function to fetch JSON data from IndexedDB using a key
// export async function fetchJSON(key) {
//   const db = await openDB();
//   return new Promise((resolve, reject) => {
//     const transaction = db.transaction(STORE_NAME, "readonly");
//     const store = transaction.objectStore(STORE_NAME);
//     const request = store.get(key);

//     request.onsuccess = () => resolve(request.result || null);
//     request.onerror = () => reject(request.error);
//   });
// }
