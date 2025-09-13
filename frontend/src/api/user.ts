import type { NewUser, User } from "../types/user";

export const fetchUsers = async (): Promise<User[]> => {
  const response = await fetch('https://jsonplaceholder.typicode.com/users');
  if (!response.ok) {
    throw new Error('Failed to fetch users');
  }
  return response.json();
};

export const createUser = async (user: NewUser): Promise<User> => {
  const response = await fetch('https://jsonplaceholder.typicode.com/users', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      ...user,
      id: Date.now(),
      company: user.company || { name: 'New Company' }
    }),
  });
  if (!response.ok) {
    throw new Error('Failed to create user');
  }
  return response.json();
};

