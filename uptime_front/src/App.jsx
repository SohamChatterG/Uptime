import { Routes, Route } from 'react-router-dom'; // <-- 1. Import routing components
import { AuthProvider, useAuth } from './context/AuthContext';
import AuthForm from './components/auth/AuthForm';
import Dashboard from './components/dashboard/Dashboard';
import AuthCallback from './components/auth/AuthCallback'; // <-- 2. Import the new callback component
import { Spinner } from './components/ui/Spinner';

function App() {
  return (
    <AuthProvider>
      <MainApp />
    </AuthProvider>
  );
}

const MainApp = () => {
  const { user, loading } = useAuth();

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <Spinner />
      </div>
    );
  }
  return (
    <Routes>
      <Route path="/" element={user ? <Dashboard /> : <AuthForm />} />
      <Route path="/auth/callback" element={<AuthCallback />} />
    </Routes>
  );
};

export default App;