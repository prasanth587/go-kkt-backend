import AsyncStorage from '@react-native-async-storage/async-storage';
import { Slot, useRouter, useSegments } from 'expo-router';
import { useEffect, useState } from 'react';
import { ActivityIndicator, View } from 'react-native';

const PRIMARY = '#035284';

export default function RootLayout() {
  const router = useRouter();
  const segments = useSegments();
  const [checkingAuth, setCheckingAuth] = useState(true);

  useEffect(() => {
    (async () => {
      try {
        const storedUser = await AsyncStorage.getItem('userData');
        const isLoggedIn = !!storedUser;
        const firstSegment = segments[0] || '';

        if (isLoggedIn) {
          // ✅ If user is logged in but on login screen → redirect to /home
          if (firstSegment === '' || firstSegment === 'login') {
            router.replace('/home');
          }
        } else {
          // ✅ If not logged in and not already on /login → redirect to login
          if (firstSegment !== 'login') {
            router.replace('/login');
          }
        }
      } catch (error) {
        console.error('Auth check error:', error);
        router.replace('/login');
      } finally {
        setCheckingAuth(false);
      }
    })();
  }, [segments]);

  if (checkingAuth) {
    return (
      <View
        style={{
          flex: 1,
          justifyContent: 'center',
          alignItems: 'center',
          backgroundColor: '#fff',
        }}
      >
        <ActivityIndicator size="large" color={PRIMARY} />
      </View>
    );
  }

  return <Slot />;
}
