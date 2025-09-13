import { useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { createUser, fetchUsers } from '../api/user';
import { useDarkMode } from '../hooks/UseDarkMode';
import type { NewUser, ToastState, User } from '../types/user';
import AddUserModal from '../components/AddUserModal';
import Toast from '../components/Toast';

export default function UserManagement() {
  const [searchTerm, setSearchTerm] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [usersPerPage] = useState(5);
  const [isDarkMode, setIsDarkMode] = useDarkMode();
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [toast, setToast] = useState<ToastState>({
    message: '',
    type: 'success',
    isVisible: false
  });

  const queryClient = useQueryClient();
  
  const { data: users = [], isLoading, error, isError } = useQuery({
    queryKey: ['users'],
    queryFn: fetchUsers,
  });

    
  const createUserMutation = useMutation({
    mutationFn: createUser,
    onSuccess: (newUser) => {
      queryClient.setQueryData(['users'], (old: User[] = []) => [...old, newUser]);
      setToast({ message: 'User added successfully!', type: 'success', isVisible: true });
    },
    onError: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      setToast({ message: 'Failed to add user. Please try again.', type: 'error', isVisible: true });
    },
  });


  const filteredUsers = users.filter(user =>
    user.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    user.email.toLowerCase().includes(searchTerm.toLowerCase()) ||
    user.company.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const indexOfLastUser = currentPage * usersPerPage;
  const indexOfFirstUser = indexOfLastUser - usersPerPage;
  const currentUsers = filteredUsers.slice(indexOfFirstUser, indexOfLastUser);
  const totalPages = Math.ceil(filteredUsers.length / usersPerPage);

  const handleAddUser = (user: NewUser) => {
    createUserMutation.mutate(user);
  };

  const paginate = (pageNumber: number) => {
    setCurrentPage(pageNumber);
  };

  if (isError) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800 flex items-center justify-center p-4">
        <div className="bg-red-50 dark:bg-red-900/20 border-2 border-red-200 dark:border-red-800 text-red-700 dark:text-red-300 px-8 py-6 rounded-2xl shadow-lg max-w-md w-full">
          <div className="flex items-center space-x-3">
            <svg className="h-8 w-8 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <div>
              <h3 className="text-lg font-semibold">Error Loading Data</h3>
              <p className="text-sm opacity-90">{error instanceof Error ? error.message : 'Failed to load users'}</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 via-blue-50 to-purple-50 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900 transition-colors duration-300">
      <div className="container mx-auto px-4 py-8">
        
      
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-4xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent mb-2">
              User Management
            </h1>
            <p className="text-gray-600 dark:text-gray-400">Manage your users with style and efficiency</p>
          </div>
          <button
            onClick={() => setIsDarkMode(!isDarkMode)}
            className="p-3 bg-white dark:bg-gray-800 rounded-full shadow-lg hover:shadow-xl transform hover:scale-105 transition-all duration-200 border border-gray-200 dark:border-gray-700"
          >
            {isDarkMode ? (
              <svg className="h-6 w-6 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
            ) : (
              <svg className="h-6 w-6 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
              </svg>
            )}
          </button>
        </div>
        
        <div className="flex flex-col lg:flex-row gap-6 mb-8">
          <div className="relative flex-1 max-w-md">
            <input
              type="text"
              placeholder="Search users, emails, companies..."
              value={searchTerm}
              onChange={(e) => {
                setSearchTerm(e.target.value);
                setCurrentPage(1);
              }}
              className="w-full pl-12 pr-4 py-4 bg-white dark:bg-gray-800 border-2 border-gray-200 dark:border-gray-700 rounded-2xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white transition-all duration-200 shadow-lg hover:shadow-xl"
            />
            <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
              <svg className="h-6 w-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            {searchTerm && (
              <button
                onClick={() => {
                  setSearchTerm('');
                  setCurrentPage(1);
                }}
                className="absolute inset-y-0 right-0 pr-4 flex items-center"
              >
                <svg className="h-5 w-5 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            )}
          </div>
          <button
            onClick={() => setIsAddModalOpen(true)}
            disabled={createUserMutation.isPending}
            className="bg-gradient-to-r from-blue-600 to-purple-600 text-white px-8 py-4 rounded-2xl hover:from-blue-700 hover:to-purple-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-3 transform hover:scale-105 transition-all duration-200 shadow-lg hover:shadow-xl font-semibold"
          >
            <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            Add New User
          </button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <section className="bg-white dark:bg-gray-800 p-6 rounded-2xl shadow-lg border border-gray-200 dark:border-gray-700">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 dark:text-gray-400 text-sm font-medium">Total Users</p>
                <p className="text-3xl font-bold text-gray-900 dark:text-white">{users.length}</p>
              </div>
              <div className="bg-blue-100 dark:bg-blue-900/30 p-3 rounded-xl">
                <svg className="h-8 w-8 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
              </div>
            </div>
          </section>
          
          <section className="bg-white dark:bg-gray-800 p-6 rounded-2xl shadow-lg border border-gray-200 dark:border-gray-700">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 dark:text-gray-400 text-sm font-medium">Filtered Results</p>
                <p className="text-3xl font-bold text-gray-900 dark:text-white">{filteredUsers.length}</p>
              </div>
              <div className="bg-green-100 dark:bg-green-900/30 p-3 rounded-xl">
                <svg className="h-8 w-8 text-green-600 dark:text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
            </div>
          </section>
          
          <section className="bg-white dark:bg-gray-800 p-6 rounded-2xl shadow-lg border border-gray-200 dark:border-gray-700">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 dark:text-gray-400 text-sm font-medium">Current Page</p>
                <p className="text-3xl font-bold text-gray-900 dark:text-white">{currentPage} / {totalPages || 1}</p>
              </div>
              <div className="bg-purple-100 dark:bg-purple-900/30 p-3 rounded-xl">
                <svg className="h-8 w-8 text-purple-600 dark:text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
              </div>
            </div>
          </section>
        </div>

        <div className="bg-white dark:bg-gray-800 shadow-2xl rounded-2xl overflow-hidden border border-gray-200 dark:border-gray-700">
          <div className="overflow-x-auto">
            <table className="min-w-full">
              <thead className="bg-gradient-to-r from-gray-50 to-gray-100 dark:from-gray-700 dark:to-gray-800">
                <tr>
                  <th className="px-8 py-6 text-left text-xs font-bold text-gray-600 dark:text-gray-300 uppercase tracking-wider">ID</th>
                  <th className="px-8 py-6 text-left text-xs font-bold text-gray-600 dark:text-gray-300 uppercase tracking-wider">Name</th>
                  <th className="px-8 py-6 text-left text-xs font-bold text-gray-600 dark:text-gray-300 uppercase tracking-wider">Email</th>
                  <th className="px-8 py-6 text-left text-xs font-bold text-gray-600 dark:text-gray-300 uppercase tracking-wider">Company</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                {isLoading ? (
                  <tr>
                    <td colSpan={5} className="px-8 py-16 text-center">
                      <div className="flex flex-col items-center justify-center space-y-4">
                        <div className="relative">
                          <div className="animate-spin rounded-full h-16 w-16 border-4 border-gray-300 dark:border-gray-600"></div>
                          <div className="animate-spin rounded-full h-16 w-16 border-4 border-blue-600 border-t-transparent absolute top-0"></div>
                        </div>
                        <p className="text-lg text-gray-600 dark:text-gray-400 font-medium">Loading amazing users...</p>
                      </div>
                    </td>
                  </tr>
                ) : currentUsers.length === 0 ? (
                  <tr>
                    <td colSpan={5} className="px-8 py-16 text-center">
                      <div className="flex flex-col items-center space-y-4">
                        <svg className="h-16 w-16 text-gray-400 dark:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                        </svg>
                        <div className="text-center">
                          <p className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
                            {searchTerm ? 'No users found' : 'No users available'}
                          </p>
                          <p className="text-gray-600 dark:text-gray-400">
                            {searchTerm ? 'Try adjusting your search terms.' : 'Add some users to get started.'}
                          </p>
                        </div>
                      </div>
                    </td>
                  </tr>
                ) : (
                  currentUsers.map((user, index) => (
                    <tr 
                      key={`${user.id}-${user.email}-${index}`}
                      className="hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors duration-200 group"
                      style={{
                        animationDelay: `${index * 50}ms`,
                        animation: 'fadeInUp 0.5s ease-out forwards'
                      }}
                    >
                      <td className="px-8 py-6 whitespace-nowrap">
                        <span className="inline-flex items-center justify-center w-10 h-10 bg-gradient-to-br from-blue-500 to-purple-600 text-white rounded-full font-bold text-sm">
                          {user.id}
                        </span>
                      </td>
                      <td className="px-8 py-6 whitespace-nowrap">
                        <div className="flex items-center">
                          <div className="w-12 h-12 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-full flex items-center justify-center text-white font-bold text-lg mr-4">
                            {user.name.charAt(0)}
                          </div>
                          <div>
                            <p className="text-lg font-semibold text-gray-900 dark:text-white group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors duration-200">
                              {user.name}
                            </p>
                          </div>
                        </div>
                      </td>
                      <td className="px-8 py-6 whitespace-nowrap">
                        <div className="flex items-center">
                          <svg className="h-5 w-5 text-gray-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207" />
                          </svg>
                          <span className="text-gray-900 dark:text-white">{user.email}</span>
                        </div>
                      </td>
                      <td className="px-8 py-6 whitespace-nowrap">
                        <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-200">
                          {user.company.name}
                        </span>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>

          {filteredUsers.length > usersPerPage && (
            <div className="px-8 py-6 bg-gradient-to-r from-gray-50 via-white to-gray-50 dark:from-gray-700 dark:via-gray-800 dark:to-gray-700 border-t border-gray-200 dark:border-gray-600">
              <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
                <div className="flex items-center space-x-4">
                  <p className="text-sm font-medium text-gray-700 dark:text-gray-300">
                    Showing <span className="font-bold text-blue-600 dark:text-blue-400">{indexOfFirstUser + 1}</span> to{' '}
                    <span className="font-bold text-blue-600 dark:text-blue-400">{Math.min(indexOfLastUser, filteredUsers.length)}</span> of{' '}
                    <span className="font-bold text-blue-600 dark:text-blue-400">{filteredUsers.length}</span> results
                  </p>
                </div>
                
                <div className="flex items-center space-x-1">
                  <button
                    onClick={() => paginate(1)}
                    disabled={currentPage === 1}
                    className="p-2 rounded-xl bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-600 hover:bg-blue-50 dark:hover:bg-blue-900/30 hover:border-blue-300 dark:hover:border-blue-600 disabled:opacity-40 disabled:cursor-not-allowed transition-all duration-200 group"
                    title="First page"
                  >
                    <svg className="h-4 w-4 text-gray-600 dark:text-gray-400 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 19l-7-7 7-7m8 14l-7-7 7-7" />
                    </svg>
                  </button>

                  <button
                    onClick={() => paginate(currentPage - 1)}
                    disabled={currentPage === 1}
                    className="px-4 py-2 rounded-xl bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-600 hover:bg-blue-50 dark:hover:bg-blue-900/30 hover:border-blue-300 dark:hover:border-blue-600 disabled:opacity-40 disabled:cursor-not-allowed transition-all duration-200 group flex items-center space-x-2"
                  >
                    <svg className="h-4 w-4 text-gray-600 dark:text-gray-400 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                    </svg>
                    <span className="text-sm font-medium text-gray-700 dark:text-gray-300 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">Previous</span>
                  </button>

                  <div className="flex items-center space-x-1">
                    {(() => {
                      const pages = [];
                      const showPages = 5;
                      let startPage = Math.max(1, currentPage - Math.floor(showPages / 2));
                      const endPage = Math.min(totalPages, startPage + showPages - 1);
                      
                      if (endPage - startPage + 1 < showPages) {
                        startPage = Math.max(1, endPage - showPages + 1);
                      }

                      if (startPage > 1) {
                        pages.push(
                          <button
                            key={1}
                            onClick={() => paginate(1)}
                            className="w-10 h-10 rounded-xl bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-600 hover:bg-blue-50 dark:hover:bg-blue-900/30 hover:border-blue-300 dark:hover:border-blue-600 text-sm font-medium text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 transition-all duration-200"
                          >
                            1
                          </button>
                        );
                        if (startPage > 2) {
                          pages.push(
                            <span key="ellipsis1" className="px-2 py-2 text-gray-400 dark:text-gray-500 text-sm">
                              ...
                            </span>
                          );
                        }
                      }

                      for (let i = startPage; i <= endPage; i++) {
                        pages.push(
                          <button
                            key={i}
                            onClick={() => paginate(i)}
                            className={`w-10 h-10 rounded-xl text-sm font-medium transition-all duration-200 transform hover:scale-105 ${
                              i === currentPage
                                ? 'bg-gradient-to-r from-blue-600 to-purple-600 text-white shadow-lg shadow-blue-500/25 dark:shadow-blue-500/20 border-0'
                                : 'bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-600 hover:bg-blue-50 dark:hover:bg-blue-900/30 hover:border-blue-300 dark:hover:border-blue-600 text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400'
                            }`}
                          >
                            {i}
                          </button>
                        );
                      }

                      if (endPage < totalPages) {
                        if (endPage < totalPages - 1) {
                          pages.push(
                            <span key="ellipsis2" className="px-2 py-2 text-gray-400 dark:text-gray-500 text-sm">
                              ...
                            </span>
                          );
                        }
                        pages.push(
                          <button
                            key={totalPages}
                            onClick={() => paginate(totalPages)}
                            className="w-10 h-10 rounded-xl bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-600 hover:bg-blue-50 dark:hover:bg-blue-900/30 hover:border-blue-300 dark:hover:border-blue-600 text-sm font-medium text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 transition-all duration-200"
                          >
                            {totalPages}
                          </button>
                        );
                      }

                      return pages;
                    })()}
                  </div>

                  <button
                    onClick={() => paginate(currentPage + 1)}
                    disabled={currentPage === totalPages}
                    className="px-4 py-2 rounded-xl bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-600 hover:bg-blue-50 dark:hover:bg-blue-900/30 hover:border-blue-300 dark:hover:border-blue-600 disabled:opacity-40 disabled:cursor-not-allowed transition-all duration-200 group flex items-center space-x-2"
                  >
                    <span className="text-sm font-medium text-gray-700 dark:text-gray-300 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">Next</span>
                    <svg className="h-4 w-4 text-gray-600 dark:text-gray-400 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                    </svg>
                  </button>

                  <button
                    onClick={() => paginate(totalPages)}
                    disabled={currentPage === totalPages}
                    className="p-2 rounded-xl bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-600 hover:bg-blue-50 dark:hover:bg-blue-900/30 hover:border-blue-300 dark:hover:border-blue-600 disabled:opacity-40 disabled:cursor-not-allowed transition-all duration-200 group"
                    title="Last page"
                  >
                    <svg className="h-4 w-4 text-gray-600 dark:text-gray-400 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 5l7 7-7 7M5 5l7 7-7 7" />
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
      <AddUserModal
        isOpen={isAddModalOpen} 
        onClose={() => setIsAddModalOpen(false)} 
        onSubmit={handleAddUser} 
      />

      <Toast
        message={toast.message} 
        type={toast.type} 
        isVisible={toast.isVisible} 
        onClose={() => setToast(prev => ({ ...prev, isVisible: false }))} 
      />
    </div>
  );
};

