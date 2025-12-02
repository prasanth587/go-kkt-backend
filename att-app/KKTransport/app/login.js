    import AsyncStorage from '@react-native-async-storage/async-storage';
import { useRouter } from 'expo-router';
import LottieView from 'lottie-react-native';
import { Eye, EyeOff, KeyRound, Lock, Mail } from 'lucide-react-native';
import { useEffect, useRef, useState } from 'react';
import {
  Animated,
  Dimensions,
  Image,
  Modal,
  Platform,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';
import { KeyboardAwareScrollView } from 'react-native-keyboard-aware-scroll-view';
import kkLogo from '../assets/app_logo1.png';
import loginAnim from '../assets/login.json';
import LoadingOverlay from '../components/LoadingOverlay.js';

const PRIMARY = '#035284';
const DEFAULT_HOST = 'http://192.168.1.45:9005/';
const { height: SCREEN_HEIGHT } = Dimensions.get('window');

export default function Login() {
  const router = useRouter();
  const [emailOrMobile, setEmailOrMobile] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const [tapCount, setTapCount] = useState(0);
  const [devModal, setDevModal] = useState(false);
  const [forgotModal, setForgotModal] = useState(false);
  const [baseUrl, setBaseUrl] = useState(DEFAULT_HOST);
  const [tempHost, setTempHost] = useState('');

  const [forgotLogin, setForgotLogin] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [showNewPassword, setShowNewPassword] = useState(false);

  // Animated popup
  const popupAnim = useRef(new Animated.Value(0)).current;
  const [popup, setPopup] = useState({ visible: false, type: 'success', message: '' });

  const showPopup = (type, message, duration = 2500) => {
    setPopup({ visible: true, type, message });
    Animated.spring(popupAnim, { toValue: 1, useNativeDriver: true }).start();
    setTimeout(() => {
      Animated.spring(popupAnim, { toValue: 0, useNativeDriver: true }).start(() =>
        setPopup({ visible: false, type: '', message: '' })
      );
    }, duration);
  };

  // Load saved base URL
  useEffect(() => {
    (async () => {
      const stored = await AsyncStorage.getItem('baseUrl');
      if (stored) setBaseUrl(stored);
    })();
  }, []);

  // ===== LOGIN =====
  const handleLogin = async () => {
    if (!emailOrMobile || !password) {
      showPopup('error', 'Please enter both email/mobile and password.');
      return;
    }

    try {
      setIsLoading(true);
      const response = await fetch(`${baseUrl}v1/adminLogin`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email_id: emailOrMobile, password }),
      });
      const data = await response.json();

      if (response.ok && data.message?.includes('User logged in successfully')) {
        await AsyncStorage.setItem('userData', JSON.stringify(data));
        await AsyncStorage.setItem('baseUrl', baseUrl);
        showPopup('success', 'Login successful!');
        setTimeout(() => router.replace('/home'), 1000);
      } else {
        showPopup('error', data.message || 'Invalid credentials.');
      }
    } catch (error) {
      showPopup('error', 'Unable to connect to server.');
    } finally {
      setIsLoading(false);
    }
  };

  // ===== 6 TAP DEVELOPER TRIGGER =====
  const handleBrandTap = () => {
    setTapCount((prev) => {
      const next = prev + 1;
      if (next >= 6) {
        setTapCount(0);
        setDevModal(true);
      }
      return next;
    });
  };

  // ===== SAVE HOST =====
  const saveNewHost = async () => {
    let trimmed = tempHost.trim();
    if (!trimmed) return showPopup('error', 'Please enter a valid host IP');

    // ✅ Extract only IP if user enters full URL
    const ipMatch = trimmed.match(/(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})/);
    const ip = ipMatch ? ipMatch[1] : trimmed;

    // ✅ Construct full base URL
    const finalBase = `http://${ip}:9005/`;

    setBaseUrl(finalBase);
    await AsyncStorage.setItem('baseUrl', finalBase);
    setTempHost('');
    setDevModal(false);
    showPopup('success', `Base URL set to ${finalBase}`);
  };


  // ===== FORGOT PASSWORD =====
  const handleForgotPassword = async () => {
    if (!forgotLogin || !newPassword) {
      showPopup('error', 'Please fill in both fields.');
      return;
    }

    try {
      setIsLoading(true);
      const response = await fetch(`${baseUrl}forgot/v1/password`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_login: forgotLogin,
          new_password: newPassword,
        }),
      });

      const data = await response.json();

      if (response.ok) {
        setForgotModal(false);
        showPopup('success', data.message || 'Password reset successful.');
        setTimeout(() => {
          showPopup('success', 'Password reset successful — go back to login?', 3000);
        }, 800);
      } else {
        showPopup('error', data.message || 'Failed to reset password.');
      }
    } catch (error) {
      showPopup('error', 'Unable to reach server.');
    } finally {
      setIsLoading(false);
    }
  };

  // ===== HEALTH CHECK =====
  const handleHealthCheck = async () => {
    try {
      const fullUrl = `${baseUrl}ping`;
      setIsLoading(true);
      const response = await fetch(fullUrl);
      const text = await response.text();
      showPopup('success', `✅ Health Check Passed\nURL: ${fullUrl}\nResponse: ${text}`);
    } catch (err) {
      showPopup('error', `❌ Health Check Failed\nURL: ${baseUrl}ping\nError: ${err.message}`);
    } finally {
      setIsLoading(false);
    }
  };


  return (
    <View style={styles.container}>
      <KeyboardAwareScrollView
        contentContainerStyle={[styles.scrollContent, { minHeight: SCREEN_HEIGHT * 0.95 }]}
        enableOnAndroid
        extraScrollHeight={Platform.OS === 'ios' ? 60 : 100}
        showsVerticalScrollIndicator={false}
      >
        {/* ===== HEADER ===== */}
        <View style={styles.topSection}>
          <TouchableOpacity onPress={handleBrandTap} activeOpacity={0.8}>
            <View style={styles.brandRow}>
              <Image source={kkLogo} style={styles.brandLogo} resizeMode="contain" />
              <Text style={styles.brandTitle}>TRANSPORT</Text>
            </View>
          </TouchableOpacity>
          <LottieView source={loginAnim} autoPlay loop style={styles.anim} />
        </View>

        {/* ===== FORM ===== */}
        <View style={styles.formSection}>
          <Text style={styles.title}>Welcome Back !</Text>
          <Text style={styles.subtitle}>Enter your Email / Mobile and Password</Text>

          {/* Email/Mobile Input */}
          <View style={styles.inputBox}>
            <Mail size={20} color={PRIMARY} style={styles.inputIcon} />
            <TextInput
              style={styles.input}
              placeholder="Email or Mobile Number"
              value={emailOrMobile}
              onChangeText={setEmailOrMobile}
              autoCapitalize="none"
            />
          </View>

          {/* Password */}
          <View style={styles.inputBox}>
            <Lock size={20} color={PRIMARY} style={styles.inputIcon} />
            <TextInput
              style={styles.input}
              placeholder="Password"
              secureTextEntry={!showPassword}
              value={password}
              onChangeText={setPassword}
            />
            <TouchableOpacity onPress={() => setShowPassword(!showPassword)}>
              {showPassword ? (
                <EyeOff size={20} color={PRIMARY} />
              ) : (
                <Eye size={20} color={PRIMARY} />
              )}
            </TouchableOpacity>
          </View>

          <TouchableOpacity
            style={[styles.loginBtn, isLoading && { opacity: 0.7 }]}
            onPress={handleLogin}
            disabled={isLoading}
          >
            <Text style={styles.loginText}>{isLoading ? 'Logging in...' : 'LOG IN'}</Text>
          </TouchableOpacity>

          <TouchableOpacity onPress={() => setForgotModal(true)}>
            <Text style={styles.forgot}>Forgot password?</Text>
          </TouchableOpacity>
        </View>
      </KeyboardAwareScrollView>

      {/* ===== DEVELOPER HOST MODAL ===== */}
      <Modal visible={devModal} transparent animationType="fade" onRequestClose={() => setDevModal(false)}>
        <View style={styles.modalOverlay}>
          <View style={styles.modalBox}>
            <Text style={styles.modalTitle}>Set API Host</Text>
            <Text style={styles.modalSub}>Current: {baseUrl || DEFAULT_HOST}</Text>

            <TextInput
              style={styles.modalInput}
              placeholder="Enter IP address only (e.g. 192.168.1.8)"
              value={tempHost}
              onChangeText={setTempHost}
              autoCapitalize="none"
            />

            <Text style={styles.exampleText}>Example: 192.168.1.8 → becomes http://192.168.1.8:9005/</Text>

            <View style={[styles.modalActions, { justifyContent: 'space-between' }]}>
              <TouchableOpacity
                onPress={handleHealthCheck}
                style={[styles.modalBtn, { backgroundColor: '#009688', flex: 1, marginRight: 8 }]}
              >
                <Text style={[styles.modalBtnText, { color: '#fff', textAlign: 'center' }]}>Health Check</Text>
              </TouchableOpacity>

              <TouchableOpacity
                onPress={() => setDevModal(false)}
                style={[styles.modalBtn, { backgroundColor: '#ccc', flex: 1, marginRight: 8 }]}
              >
                <Text style={[styles.modalBtnText, { textAlign: 'center' }]}>Cancel</Text>
              </TouchableOpacity>

              <TouchableOpacity
                onPress={saveNewHost}
                style={[styles.modalBtn, { backgroundColor: PRIMARY, flex: 1 }]}
              >
                <Text style={[styles.modalBtnText, { color: '#fff', textAlign: 'center' }]}>Save</Text>
              </TouchableOpacity>
            </View>

          </View>
        </View>
      </Modal>

      {/* ===== FORGOT PASSWORD MODAL ===== */}
      <Modal visible={forgotModal} transparent animationType="fade" onRequestClose={() => setForgotModal(false)}>
        <View style={styles.modalOverlay}>
          <View style={styles.modalBox}>
            <Text style={styles.modalTitle}>Reset Password</Text>
            <Text style={styles.modalSub}>Enter your Mobile/Email and New Password</Text>

            {/* Login Field */}
            <View style={styles.inputBox}>
              <KeyRound size={20} color={PRIMARY} style={styles.inputIcon} />
              <TextInput
                style={styles.input}
                placeholder="Email or Mobile Number"
                value={forgotLogin}
                onChangeText={setForgotLogin}
                autoCapitalize="none"
              />
            </View>

            {/* New Password Field */}
            <View style={styles.inputBox}>
              <Lock size={20} color={PRIMARY} style={styles.inputIcon} />
              <TextInput
                style={styles.input}
                placeholder="New Password"
                secureTextEntry={!showNewPassword}
                value={newPassword}
                onChangeText={setNewPassword}
              />
              <TouchableOpacity onPress={() => setShowNewPassword(!showNewPassword)}>
                {showNewPassword ? (
                  <EyeOff size={20} color={PRIMARY} />
                ) : (
                  <Eye size={20} color={PRIMARY} />
                )}
              </TouchableOpacity>
            </View>

            <View style={styles.modalActions}>
              <TouchableOpacity onPress={() => setForgotModal(false)} style={[styles.modalBtn, { backgroundColor: '#ccc' }]}>
                <Text style={styles.modalBtnText}>Cancel</Text>
              </TouchableOpacity>
              <TouchableOpacity onPress={handleForgotPassword} style={[styles.modalBtn, { backgroundColor: PRIMARY }]}>
                <Text style={[styles.modalBtnText, { color: '#fff' }]}>Submit</Text>
              </TouchableOpacity>
            </View>
          </View>
        </View>
      </Modal>

      {/* ===== POPUP ===== */}
      {popup.visible && (
        <Animated.View
          style={[
            styles.popup,
            {
              borderColor: popup.type === 'success' ? PRIMARY : '#d32f2f',
              transform: [{ scale: popupAnim }],
            },
          ]}
        >
          <Text style={[styles.popupText, { color: popup.type === 'success' ? PRIMARY : '#d32f2f' }]}>
            {popup.message}
          </Text>
        </Animated.View>
      )}

      <LoadingOverlay visible={isLoading} />
    </View>
  );
}

// ===== STYLES =====
const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#fff' },
  scrollContent: { flexGrow: 1, justifyContent: 'center', paddingHorizontal: 40, paddingVertical: 30 },
  topSection: { alignItems: 'center', marginBottom: 40 },
  brandRow: { flexDirection: 'row', alignItems: 'center', justifyContent: 'center', marginBottom: 15 },
  brandLogo: { width: 60, height: 60, marginRight: 10 },
  brandTitle: { fontSize: 26, fontWeight: '800', color: PRIMARY },
  anim: { width: 120, height: 120, marginTop: 10 },
  formSection: { alignItems: 'center', width: '100%' },
  title: { fontSize: 26, fontWeight: '700', color: PRIMARY, marginBottom: 4 },
  subtitle: { fontSize: 15, color: '#666', marginBottom: 25 },
  inputBox: {
    flexDirection: 'row',
    alignItems: 'center',
    borderColor: '#ddd',
    borderWidth: 1,
    borderRadius: 12,
    backgroundColor: '#f8f8f8',
    marginBottom: 18,
    width: '100%',
    paddingHorizontal: 10,
  },
  inputIcon: { marginRight: 8 },
  input: { flex: 1, fontSize: 16, paddingVertical: 12 },
  loginBtn: { width: '100%', backgroundColor: PRIMARY, paddingVertical: 15, borderRadius: 12, alignItems: 'center', marginTop: 10 },
  loginText: { color: '#fff', fontWeight: 'bold', fontSize: 16 },
  forgot: { color: PRIMARY, textDecorationLine: 'underline', marginTop: 15, fontSize: 14, fontWeight: '600' },
  modalOverlay: { flex: 1, backgroundColor: 'rgba(0,0,0,0.45)', justifyContent: 'center', alignItems: 'center' },
  modalBox: { width: '85%', backgroundColor: '#fff', borderRadius: 15, padding: 20, elevation: 6 },
  modalTitle: { fontSize: 18, fontWeight: '700', color: PRIMARY },
  modalSub: { fontSize: 13, color: '#666', marginVertical: 8 },
  modalInput: { borderWidth: 1, borderColor: '#ccc', borderRadius: 8, padding: 10, fontSize: 13, marginBottom: 15 },
  modalActions: { flexDirection: 'row', justifyContent: 'flex-end', gap: 10 },
  modalBtn: { borderRadius: 8, paddingVertical: 10, paddingHorizontal: 18 },
  modalBtnText: { fontWeight: '700', fontSize: 15, color: '#000' },
  popup: {
    position: 'absolute',
    bottom: 60,
    left: '10%',
    width: '80%',
    backgroundColor: '#fff',
    borderWidth: 2,
    borderRadius: 14,
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: 15,
    paddingHorizontal: 20,
    elevation: 8,
  },
  popupText: { fontSize: 15, fontWeight: '700', textAlign: 'center' },
  exampleText: {
    fontSize: 10,
    fontWeight: '700',
    color: '#777',
    marginTop: -8,
    marginBottom: 15,
    fontStyle: 'italic',
  },
});
