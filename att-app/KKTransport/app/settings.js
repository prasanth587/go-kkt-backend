import AsyncStorage from '@react-native-async-storage/async-storage';
import { useRouter } from 'expo-router';
import { ArrowLeft, LogOut, User } from 'lucide-react-native';
import { useEffect, useState } from 'react';
import {
  ScrollView,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
} from 'react-native';
import LoadingOverlay from '../components/LoadingOverlay';
const PRIMARY = '#035284';

export default function Settings() {
  const router = useRouter();
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const verifyAuth = async () => {
      const data = await AsyncStorage.getItem('userData');
      if (!data) router.replace('/login');
    };
    verifyAuth();
  }, []);

  useEffect(() => {
    const loadUser = async () => {
      try {
        const data = await AsyncStorage.getItem('userData');
        if (data) setUser(JSON.parse(data));
        else router.replace('/login');
      } catch (error) {
        console.error(error);
      } finally {
        setLoading(false);
      }
    };
    loadUser();
  }, []);

  if (loading) return <LoadingOverlay visible />;

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <TouchableOpacity onPress={() => router.back()}>
          <ArrowLeft size={24} color="#fff" />
        </TouchableOpacity>

        <Text style={styles.headerTitle}>Settings</Text>
        <View style={{ width: 24 }} />
      </View>

      <ScrollView contentContainerStyle={styles.scrollContent}>
        {user && (
          <View style={styles.userCard}>
            <View style={styles.avatarContainer}>
              <View style={styles.avatarCircle}>
                <Text style={styles.avatarText}>
                  {String(user?.first_name?.[0] || '?')}
                </Text>
              </View>
            </View>
            <View>
              <Text style={styles.userName}>
                {String(user?.first_name || '')} {String(user?.last_name || '')}
              </Text>
              <Text style={styles.userRole}>{String(user?.role_name || '')}</Text>
            </View>
          </View>
        )}

        <View style={styles.menuContainer}>
          <TouchableOpacity
            style={styles.menuCard}
            activeOpacity={0.85}
            onPress={() => router.push('/profile')}
          >
            <View style={styles.iconBox}>
              <User size={22} color={PRIMARY} />
            </View>
            <View style={styles.menuTextBox}>
              <Text style={styles.menuTitle}>Profile</Text>
              <Text style={styles.menuSubtitle}>View personal & company info</Text>
            </View>
          </TouchableOpacity>

          <TouchableOpacity
            style={styles.menuCard}
            activeOpacity={0.85}
            onPress={() => router.push('/logout')}
          >
            <View style={[styles.iconBox, { backgroundColor: '#fdecea' }]}>
              <LogOut size={22} color="#d32f2f" />
            </View>
            <View style={styles.menuTextBox}>
              <Text style={[styles.menuTitle, { color: '#d32f2f' }]}>
                Logout
              </Text>
              <Text style={styles.menuSubtitle}>
                Sign out and end your session
              </Text>
            </View>
          </TouchableOpacity>
        </View>
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#fff' },

  header: {
    backgroundColor: PRIMARY,
    paddingTop: 55,
    paddingBottom: 15,
    paddingHorizontal: 20,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  headerTitle: { fontSize: 20, fontWeight: '700', color: '#fff' },

  scrollContent: { padding: 20, paddingBottom: 50 },

  userCard: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#f5f8fb',
    borderRadius: 14,
    padding: 16,
    marginBottom: 25,
    elevation: 3,
    shadowColor: '#000',
    shadowOpacity: 0.15,
    shadowOffset: { width: 0, height: 2 },
  },
  avatarContainer: { marginRight: 15 },
  avatarCircle: {
    width: 40,
    height: 40,
    borderRadius: 30,
    backgroundColor: PRIMARY,
    justifyContent: 'center',
    alignItems: 'center',
  },
  avatarText: { color: '#fff', fontSize: 20, fontWeight: '700' },
  userName: { fontSize: 18, fontWeight: '700', color: PRIMARY },
  userRole: { fontSize: 15, color: '#555', marginTop: 2 },
  userOrg: { fontSize: 13, color: '#777', marginTop: 2 },

  menuContainer: { gap: 16 },
  menuCard: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#fff',
    padding: 16,
    borderRadius: 12,
    elevation: 3,
    shadowColor: '#000',
    shadowOpacity: 0.1,
    shadowOffset: { width: 0, height: 1 },
  },
  iconBox: {
    backgroundColor: '#e9f2f9',
    borderRadius: 10,
    width: 45,
    height: 45,
    justifyContent: 'center',
    alignItems: 'center',
  },
  menuTextBox: { marginLeft: 15 },
  menuTitle: { fontSize: 17, fontWeight: '700', color: PRIMARY },
  menuSubtitle: { fontSize: 13, color: '#666', marginTop: 2 },
});
