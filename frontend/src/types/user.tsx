// Types
export interface User {
  id: number;
  name: string;
  email: string;
  company: {
    name: string;
  };
}

export interface NewUser {
  name: string;
  email: string;
  company?: {
    name: string;
  };
}

export interface ToastState {
  message: string;
  type: 'success' | 'error';
  isVisible: boolean;
}

export interface AddUserModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (user: NewUser) => void;
}

export interface ToastProps {
  message: string;
  type: 'success' | 'error';
  isVisible: boolean;
  onClose: () => void;
}