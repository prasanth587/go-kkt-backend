import AsyncStorage from '@react-native-async-storage/async-storage';
import * as Location from 'expo-location';
import { useRouter } from 'expo-router';
import {
  CheckCircle,
  Clock,
  LogIn,
  LogOut,
  MoreVertical,
  XCircle
} from 'lucide-react-native';
import { useEffect, useRef, useState } from 'react';
import {
  Animated,
  Easing,
  Image,
  Modal,
  Pressable,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
} from 'react-native';
import AnimatedRe, {
  useAnimatedStyle,
  useSharedValue,
  withRepeat,
  withTiming,
} from 'react-native-reanimated';
import kkLogo from '../assets/app_logo.png';
import LoadingOverlay from '../components/LoadingOverlay';

const PRIMARY = '#035284';
const DEFAULT_BASE_URL = 'http://192.168.1.45:9005/';

export default function Home() {
  const router = useRouter();
  const [user, setUser] = useState(null);
  const [isCheckingIn, setIsCheckingIn] = useState(false);
  const [isCheckingOut, setIsCheckingOut] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  const [currentTime, setCurrentTime] = useState('');
  const [currentDate, setCurrentDate] = useState('');
  const [location, setLocation] = useState(null);
  const [baseUrl, setBaseUrl] = useState(DEFAULT_BASE_URL);
  const [popup, setPopup] = useState({ visible: false, type: '', message: '', sub: '' });
  const [attendance, setAttendance] = useState({
    check_in_time: null,
    check_out_time: null,
  });

  const popupScale = useRef(new Animated.Value(0)).current;
  const fadeAnim = useRef(new Animated.Value(0)).current;
  const refreshAnim = useRef(new Animated.Value(0)).current;
  const [greeting, setGreeting] = useState('');

  useEffect(() => {
    const verifyAuth = async () => {
      const data = await AsyncStorage.getItem('userData');
      if (!data) router.replace('/login');
    };
    verifyAuth();
  }, []);

  useEffect(() => {
    const interval = setInterval(() => {
      const now = new Date();
      const hours = now.getHours();

      // Set time
      setCurrentTime(
        now.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', hour12: true })
      );
      setCurrentDate(now.toDateString());

      // Determine greeting
      if (hours >= 5 && hours < 12) setGreeting('Good Morning');
      else if (hours >= 12 && hours < 16) setGreeting('Good Afternoon');
      else if (hours >= 16 && hours < 20) setGreeting('Good Evening');
      else setGreeting('Good Night');
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  // Utility: Working hours
  const calculateWorkingHours = (inTime, outTime) => {
    if (!inTime || !outTime) return '';
    const parse = (t) => {
      const [time, mer] = t.split(' ');
      let [h, m] = time.split(':').map(Number);
      if (mer === 'PM' && h < 12) h += 12;
      if (mer === 'AM' && h === 12) h = 0;
      return h * 60 + m;
    };
    const diff = parse(outTime) - parse(inTime);
    const hrs = Math.floor(diff / 60);
    const mins = diff % 60;
    return `${hrs}h ${mins}m`;
  };

  // Fade in
  useEffect(() => {
    Animated.timing(fadeAnim, { toValue: 1, duration: 600, useNativeDriver: true }).start();
  }, [attendance]);

  // Live clock
  useEffect(() => {
    const interval = setInterval(() => {
      const now = new Date();
      setCurrentTime(
        now.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', hour12: true })
      );
      setCurrentDate(now.toDateString());
    }, 1000);
    return () => clearInterval(interval);
  }, []);

  // Pulse animation for button
  const scale = useSharedValue(1);
  useEffect(() => {
    scale.value = withRepeat(withTiming(1.04, { duration: 1600 }), -1, true);
  }, []);
  const animatedScale = useAnimatedStyle(() => ({ transform: [{ scale: scale.value }] }));

  // Load user and location
  useEffect(() => {
    (async () => {
      try {
        const data = await AsyncStorage.getItem('userData');
        const storedBase = await AsyncStorage.getItem('baseUrl');
        if (storedBase) setBaseUrl(storedBase);
        if (data) {
          const parsed = JSON.parse(data);
          setUser(parsed);

          const empAtt = parsed.employee_attendance || {};
          setAttendance({
            check_in_time: empAtt.in_time || null,
            check_out_time: empAtt.out_time || null,
          });
        }
        else router.replace('/');

        const { status } = await Location.requestForegroundPermissionsAsync();
        if (status === 'granted') {
          const loc = await Location.getCurrentPositionAsync({});
          setLocation(loc.coords);
        }
      } catch (err) {
        console.error('Error loading data:', err);
      }
    })();
  }, []);

  // Popup handler
  const showPopup = (type, message, sub = '') => {
    setPopup({ visible: true, type, message, sub });
    Animated.spring(popupScale, { toValue: 1, useNativeDriver: true }).start();
    setTimeout(() => {
      Animated.spring(popupScale, { toValue: 0, useNativeDriver: true }).start(() =>
        setPopup({ visible: false, type: '', message: '', sub: '' })
      );
    }, 2200);
  };

  // ===== REFRESH ATTENDANCE =====
  const handleRefresh = async () => {
    if (!user) return;
    try {
      setRefreshing(true);
      Animated.loop(
        Animated.timing(refreshAnim, {
          toValue: 1,
          duration: 1000,
          easing: Easing.linear,
          useNativeDriver: true,
        })
      ).start();

      const response = await fetch(`${baseUrl}v1/adminLogin`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email_id: user.email_id || user.user_login,
          password: user.password || 'welcome', // fallback only for test
        }),
      });
      const data = await response.json();

      if (response.ok) {
        const empAtt = data.employee_attendance || {};
        setAttendance({
          check_in_time: empAtt.in_time || null,
          check_out_time: empAtt.out_time || null,
        });
        await AsyncStorage.setItem('userData', JSON.stringify(data));
        showPopup('success', 'Attendance refreshed successfully');
      } else {
        showPopup('error', data.message || 'Unable to refresh');
      }
    } catch (err) {
      showPopup('error', 'Connection Error', 'Unable to reach server.');
    } finally {
      setRefreshing(false);
      refreshAnim.stopAnimation();
      refreshAnim.setValue(0);
    }
  };

  // ===== Check In / Out =====
  const check = async (type) => {
    if (!user) return;
    const isIn = type === 'in';
    const setLoading = isIn ? setIsCheckingIn : setIsCheckingOut;

    try {
      setLoading(true);
      const now = new Date();
      const payload = {
        [`check_${type}_date_str`]: now.toISOString().split('T')[0],
        [`${isIn ? 'in' : 'out'}_time`]: now.toLocaleTimeString([], {
          hour: '2-digit',
          minute: '2-digit',
          hour12: true,
        }),
        [`check_${type}_latitude`]: location?.latitude || 0,
        [`check_${type}_longitude`]: location?.longitude || 0,
        [`check_${type}_by_id`]: user.login_id,
      };

      const endpoint = `${baseUrl}v1/create/attendance/${type}/${user.employee_id}`;
      const res = await fetch(endpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      const data = await res.json();
      if (res.ok) {
        if (isIn) {
          setAttendance((p) => ({ ...p, check_in_time: payload.in_time }));
          showPopup('success', 'Checked In Successfully', `Time: ${payload.in_time}`);
        } else {
          setAttendance((p) => ({ ...p, check_out_time: payload.out_time }));
          showPopup('success', 'Checked Out Successfully', `Time: ${payload.out_time}`);
        }
      } else {
        showPopup('error', 'Action Failed', data.message || 'Please try again.');
      }
    } catch (err) {
      showPopup('error', 'Connection Error', 'Unable to reach server.');
    } finally {
      setLoading(false);
    }
  };

  const totalHours = calculateWorkingHours(
    attendance.check_in_time,
    attendance.check_out_time
  );

  const rotateStyle = {
    transform: [
      {
        rotate: refreshAnim.interpolate({
          inputRange: [0, 1],
          outputRange: ['0deg', '360deg'],
        }),
      },
    ],
  };

  return (
    <View style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <View style={styles.headerRow}>
          <View style={styles.headerLeft}>
            <Image source={kkLogo} style={styles.logo} resizeMode="contain" />
            <View>
              <Text style={styles.brandName}>TRANSPORT</Text>
              <Text style={styles.welcome}>Welcome back 👋</Text>
            </View>
          </View>

          <View style={styles.headerIcons}>
            <TouchableOpacity onPress={() => router.push('/settings')}>
              <MoreVertical size={26} color="#fff" style={{ marginLeft: 10 }} />
            </TouchableOpacity>
          </View>
        </View>
      </View>

      {/* Main Panel */}
      <View style={styles.panel}>
        <View style={styles.timeCard}>
          <Clock size={30} color={PRIMARY} style={{ marginRight: 10 }} />
          <View>
            <Text style={styles.timeText}>{currentTime}</Text>
            <Text style={styles.dateText}>{currentDate}</Text>
          </View>
        </View>

        {user && (
          <View style={styles.userInfo}>
            <Text style={styles.userName}>
              {greeting}, {user.first_name}! 👋
            </Text>
            <Text style={styles.userSub}>{user.role_name}</Text>
          </View>
        )}

        {/* Status */}
        <Animated.View style={{ opacity: fadeAnim }}>
          <View style={styles.messageCard}>
            <Text style={styles.messageTitle}>Status Update</Text>
            <Text style={styles.messageText}>
              {attendance.check_in_time && !attendance.check_out_time
                ? `You are checked in at ${attendance.check_in_time}. Don't forget to check out!`
                : attendance.check_out_time
                  ? `You have successfully checked in at ${attendance.check_in_time} and checked out at ${attendance.check_out_time}, See you tomorrow!.`
                  : 'You haven’t checked in yet. Tap below to start your day!'}
            </Text>
          </View>
        </Animated.View>

        {/* Buttons */}
        <AnimatedRe.View style={[animatedScale, styles.animatedBtn]}>
          <Pressable
            style={[styles.button, styles.checkIn, attendance.check_in_time && { opacity: 0.6 }]}
            onPress={() => check('in')}
            disabled={!!attendance.check_in_time || isCheckingIn}
          >
            <LogIn color="#fff" size={20} style={{ marginRight: 8 }} />
            <Text style={styles.btnText}>
              {attendance.check_in_time ? 'Checked In' : isCheckingIn ? 'Checking In...' : 'Check In'}
            </Text>
          </Pressable>
        </AnimatedRe.View>

        <AnimatedRe.View style={[animatedScale, styles.animatedBtn]}>
          <Pressable
            style={[
              styles.button,
              styles.checkOut,
              attendance.check_out_time && { opacity: 0.6 },
            ]}
            onPress={() => check('out')}
            disabled={
              !attendance.check_in_time ||
              !!attendance.check_out_time ||
              isCheckingOut
            }
          >
            <LogOut color="#fff" size={20} style={{ marginRight: 8 }} />
            <Text style={styles.btnText}>
              {attendance.check_out_time
                ? 'Checked Out'
                : isCheckingOut
                  ? 'Checking Out...'
                  : 'Check Out'}
            </Text>
          </Pressable>
        </AnimatedRe.View>

        {/* Summary */}
        {attendance?.check_in_time && (
          <View style={styles.summaryCard}>
            <View style={styles.summaryRow}>
              <Text style={styles.summaryText}>
                Check-In: <Text style={styles.summaryValue}>{attendance.check_in_time}</Text>
              </Text>
              {attendance?.check_out_time && (
                <Text style={styles.summaryText}>
                  Check-Out: <Text style={styles.summaryValue}>{attendance.check_out_time}</Text>
                </Text>
              )}
            </View>
            {attendance?.check_out_time && (
              <Text style={[styles.summaryText, { color: '#000', marginTop: 5, textAlign: 'center' }]}>
                Total Working Hours:{' '}
                <Text style={[styles.summaryValue, { color: PRIMARY }]}>{totalHours || '--'}</Text>
              </Text>
            )}
          </View>
        )}
      </View>

      {/* Popup */}
      <Modal transparent visible={popup.visible} animationType="fade">
        <View style={styles.popupOverlay}>
          <Animated.View
            style={[
              styles.popupBox,
              {
                borderColor: popup.type === 'success' ? PRIMARY : '#d32f2f',
                transform: [{ scale: popupScale }],
              },
            ]}
          >
            {popup.type === 'success' ? (
              <CheckCircle size={50} color={PRIMARY} />
            ) : (
              <XCircle size={50} color="#d32f2f" />
            )}
            <Text
              style={[
                styles.popupText,
                { color: popup.type === 'success' ? PRIMARY : '#d32f2f' },
              ]}
            >
              {popup.message}
            </Text>
            {popup.sub ? (
              <Text style={[styles.popupSub, { color: '#333' }]}>{popup.sub}</Text>
            ) : null}
          </Animated.View>
        </View>
      </Modal>

      <LoadingOverlay visible={isCheckingIn || isCheckingOut || refreshing} />
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#fff' },
  backgroundAnim: { ...StyleSheet.absoluteFillObject, opacity: 0.1, position: 'absolute' },
  header: { height: 100, justifyContent: 'flex-end', paddingHorizontal: 15, backgroundColor: PRIMARY },
  headerRow: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginBottom: 10 },
  headerLeft: { flexDirection: 'row', alignItems: 'center' },
  headerIcons: { flexDirection: 'row', alignItems: 'center' },
  logo: { width: 42, height: 42, marginRight: 10, borderRadius: 20 },
  brandName: { fontSize: 18, fontWeight: '700', color: '#fff' },
  welcome: { color: '#cfe8fa', fontSize: 13, marginTop: 1, fontWeight: '500', fontStyle: 'italic' },
  panel: { flex: 1, alignItems: 'center', paddingTop: 30, borderTopLeftRadius: 25, borderTopRightRadius: 25 },
  timeCard: { flexDirection: 'row', backgroundColor: '#f4f7fa', borderRadius: 18, paddingVertical: 12, paddingHorizontal: 25, alignItems: 'center', elevation: 15 },
  timeText: { fontSize: 25, fontWeight: '800', color: PRIMARY, fontStyle: 'italic' },
  dateText: { fontSize: 15, color: '#777', fontWeight: '800', fontStyle: 'italic' },
  userInfo: { alignItems: 'center', marginTop: 25 },
  userName: { fontSize: 22, fontWeight: '700', color: PRIMARY, fontStyle: 'italic' },
  userSub: { fontSize: 18, color: '#666', marginTop: 4, fontWeight: '700' },
  animatedBtn: { width: '80%', marginVertical: 12 },
  button: { flexDirection: 'row', alignItems: 'center', justifyContent: 'center', borderRadius: 30, paddingVertical: 16, elevation: 5 },
  checkIn: { backgroundColor: PRIMARY },
  checkOut: { backgroundColor: '#d32f2f' },
  btnText: { color: '#fff', fontWeight: '700', fontSize: 20 },
  summaryCard: {
    marginTop: 15,
    backgroundColor: '#eef6fb',
    paddingVertical: 10,
    paddingHorizontal: 18,
    borderRadius: 12,
    width: '90%',
    alignItems: 'center',
  },
  summaryRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    width: '100%',
    marginBottom: 15,
  },
  summaryText: {
    color: PRIMARY,
    fontSize: 14,
    fontWeight: '600',
  },
  summaryValue: {
    fontWeight: '700',
    color: '#000',
  },
  popupOverlay: { flex: 1, justifyContent: 'center', alignItems: 'center', backgroundColor: 'rgba(0,0,0,0.45)' },
  popupBox: { width: '75%', backgroundColor: '#fff', borderRadius: 20, alignItems: 'center', paddingVertical: 25, paddingHorizontal: 15, elevation: 8, borderWidth: 2 },
  popupText: { fontSize: 18, fontWeight: '700', marginTop: 12, textAlign: 'center' },
  popupSub: { fontSize: 14, fontWeight: '600', marginTop: 6, textAlign: 'center' },
  messageCard: { backgroundColor: '#eaf3fa', borderLeftWidth: 4, borderLeftColor: PRIMARY, borderRadius: 10, paddingVertical: 15, paddingHorizontal: 20, width: '90%', marginTop: 35, marginBottom: 35, elevation: 10 },
  messageTitle: { fontSize: 18, fontWeight: '700', color: PRIMARY, marginBottom: 8 },
  messageText: { fontSize: 14, color: '#444', lineHeight: 20, fontStyle: 'italic' },
});
