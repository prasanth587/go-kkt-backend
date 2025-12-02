import AsyncStorage from '@react-native-async-storage/async-storage';
import { useRouter } from 'expo-router';
import {
  ArrowLeft,
  BadgeCheck,
  Briefcase,
  Building2,
  CalendarDays,
  Hash,
  Link,
  Mail,
  MapPin,
  Phone,
  Smartphone,
  User
} from 'lucide-react-native';
import { useEffect, useRef, useState } from 'react';
import {
  ScrollView,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
} from 'react-native';
import LoadingOverlay from '../components/LoadingOverlay.js';

const PRIMARY = '#035284';
const DEFAULT_HOST = 'http://192.168.1.45:9005/';
export default function Profile() {
  const router = useRouter();
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [devVisible, setDevVisible] = useState(false);
  const [baseUrl, setBaseUrl] = useState(DEFAULT_HOST);
  const tapCount = useRef(0);
  const lastTapTime = useRef(0);

  useEffect(() => {
    const verifyAuth = async () => {
      const data = await AsyncStorage.getItem('userData');
      if (!data) router.replace('/login');
    };
    verifyAuth();
  }, []);

  // ===== Load User & Base URL =====
  useEffect(() => {
    const loadData = async () => {
      try {
        const data = await AsyncStorage.getItem('userData');
        const storedBaseUrl = await AsyncStorage.getItem('baseUrl');
        if (data) setUser(JSON.parse(data));
        else router.replace('/login');
        if (storedBaseUrl) setBaseUrl(storedBaseUrl);
      } catch (error) {
        console.error(error);
      } finally {
        setLoading(false);
      }
    };
    loadData();
  }, []);

  // ===== Handle Developer Tap =====
  const handleHeaderTap = () => {
    const now = Date.now();
    if (now - lastTapTime.current < 1000) {
      tapCount.current += 1;
    } else {
      tapCount.current = 1;
    }
    lastTapTime.current = now;

    if (tapCount.current >= 6) {
      setDevVisible((prev) => !prev);
      tapCount.current = 0;
    }
  };

  if (loading) return <LoadingOverlay visible />;

  return (
    <View style={styles.container}>
      {/* ===== Header ===== */}
      <View style={styles.header}>
        <TouchableOpacity onPress={() => router.back()}>
          <ArrowLeft size={24} color="#fff" />
        </TouchableOpacity>

        {/* Tap 6x to toggle dev info */}
        <TouchableOpacity onPress={handleHeaderTap} activeOpacity={0.8}>
          <Text style={styles.headerTitle}>Profile</Text>
        </TouchableOpacity>

        <View style={{ width: 24 }}></View>
      </View>

      <ScrollView contentContainerStyle={styles.scrollContent}>
        {/* ===== User Details ===== */}
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>User Details</Text>

          <View style={styles.item}>
            <User size={20} color={PRIMARY} />
            <Text style={styles.label}>Name</Text>
            <Text style={styles.value}>{user.first_name || '—'}</Text>
          </View>

          <View style={styles.item}>
            <Mail size={20} color={PRIMARY} />
            <Text style={styles.label}>Email</Text>
            <Text style={styles.value}>{user.email_id || '—'}</Text>
          </View>

          <View style={styles.item}>
            <Smartphone size={20} color={PRIMARY} />
            <Text style={styles.label}>Mobile</Text>
            <Text style={styles.value}>{user.mobile_no || '—'}</Text>
          </View>

          <View style={styles.item}>
            <BadgeCheck size={20} color={PRIMARY} />
            <Text style={styles.label}>Role</Text>
            <Text style={styles.value}>{user.role_name || '—'}</Text>
          </View>

          <View style={styles.item}>
            <CalendarDays size={20} color={PRIMARY} />
            <Text style={styles.label}>Last Login</Text>
            <Text style={styles.value}>
              {user.last_login
                ? new Date(user.last_login).toLocaleString()
                : '—'}
            </Text>
          </View>
        </View>

        {/* ===== Organisation Details ===== */}
        {user.organisation && (
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Organisation Details</Text>

            <View style={styles.item}>
              <Building2 size={20} color={PRIMARY} />
              <Text style={styles.label}>Name</Text>
              <Text style={styles.value}>
                {user.organisation.display_name || '—'}
              </Text>
            </View>

            <View style={styles.item}>
              <Mail size={20} color={PRIMARY} />
              <Text style={styles.label}>Email</Text>
              <Text style={styles.value}>
                {user.organisation.email_id || '—'}
              </Text>
            </View>

            <View style={styles.item}>
              <Phone size={20} color={PRIMARY} />
              <Text style={styles.label}>Contact</Text>
              <Text style={styles.value}>
                {user.organisation.contact_name
                  ? `${user.organisation.contact_name} (${user.organisation.contact_no})`
                  : '—'}
              </Text>
            </View>

            <View style={styles.item}>
              <MapPin size={20} color={PRIMARY} />
              <Text style={styles.label}>City</Text>
              <Text style={styles.value}>{user.organisation.city || '—'}</Text>
            </View>
          </View>
        )}

        {/* ===== Developer Section: Base URLs ===== */}
        {devVisible && (
          <View style={[styles.section, { backgroundColor: '#f0f6fa' }]}>
            <Text style={styles.sectionTitle}>Developer Info</Text>

            <View style={styles.item}>
              <Briefcase size={20} color={PRIMARY} />
              <Text style={styles.label}>Employee ID</Text>
              <Text style={styles.value}>{user.employee_id || '—'}</Text>
            </View>

            <View style={styles.item}>
              <Hash size={20} color={PRIMARY} />
              <Text style={styles.label}>Login ID</Text>
              <Text style={styles.value}>{user.login_id || '—'}</Text>
            </View>

            <View style={styles.item}>
              <Link size={20} color={PRIMARY} />
              <Text style={styles.label}>Current Base URL</Text>
              <Text
                style={[styles.value, { flex: 2, fontSize: 12, color: '#333' }]}
              >
                {baseUrl || DEFAULT_HOST}
              </Text>
            </View>

            <View style={styles.item}>
              <Link size={20} color={PRIMARY} />
              <Text style={styles.label}>Default Base URL</Text>
              <Text
                style={[styles.value, { flex: 2, fontSize: 12, color: '#777' }]}
              >
                {DEFAULT_HOST}
              </Text>
            </View>
          </View>
        )}
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

  section: {
    backgroundColor: '#f8f9fa',
    borderRadius: 14,
    padding: 16,
    marginBottom: 20,
    elevation: 8,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '700',
    color: PRIMARY,
    marginBottom: 12,
  },
  item: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: 8,
    borderBottomWidth: 1,
    borderColor: '#eee',
  },
  label: {
    flex: 1,
    fontSize: 13,
    fontWeight: '800',
    color: '#333',
    marginLeft: 8,
  },
  value: {
    flex: 1.2,
    textAlign: 'right',
    fontSize: 13,
    color: '#555',
  },
});
