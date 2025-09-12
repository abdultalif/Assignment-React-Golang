import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import UserManagement from './pages/UserManagement';

export default function App() {
  const queryClient = new QueryClient();

  return (
    <QueryClientProvider client={queryClient}>
      <UserManagement />
    </QueryClientProvider>
  );
}
