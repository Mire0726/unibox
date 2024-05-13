// components/Login.tsx
import { signInWithGoogle } from '../firebase';  // 適切なパスを設定

const Login = () => {
  return (
    <div>
      <h1>Login</h1>
      <button onClick={signInWithGoogle}>Sign in with Google</button>
    </div>
  );
};

export default Login;
