import { initializeApp } from "firebase/app";
import { getAuth, onAuthStateChanged } from "firebase/auth";

const firebaseConfig = {
  apiKey: "AIzaSyA0TlYyTtwo5xA2rFCEVOdGxtWW2irIvRk",
  authDomain: "unibox-8f4f2.firebaseapp.com",
  projectId: "unibox-8f4f2",
  storageBucket: "unibox-8f4f2.appspot.com",
  messagingSenderId: "431775755273",
  appId: "1:431775755273:web:c216a5f243ea7475cadcd8",
  measurementId: "G-F78174DRDD",
};
const app = initializeApp(firebaseConfig);
const auth = getAuth(app); // 認証サービス

export { auth };
