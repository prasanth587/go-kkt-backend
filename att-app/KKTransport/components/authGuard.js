import AsyncStorage from '@react-native-async-storage/async-storage';
import { useRouter } from 'expo-router';
import { useEffect } from 'react';

export default function useAuthGuard() {
  const router = useRouter();
  useEffect(() => {
    (async () => {
      const data = await AsyncStorage.getItem('userData');
      if (!data) router.replace('/login');
    })();
  }, []);
}
