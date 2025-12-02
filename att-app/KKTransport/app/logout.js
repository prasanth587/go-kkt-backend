import AsyncStorage from '@react-native-async-storage/async-storage';
import { useRouter } from 'expo-router';
import { LogOut, XCircle } from 'lucide-react-native';
import { useState } from 'react';
import {
  Modal,
  Pressable,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
} from 'react-native';

export default function Logout() {
  const router = useRouter();
  const [modalVisible, setModalVisible] = useState(true);

  const confirmLogout = async () => {
    try {
      await AsyncStorage.clear();
      setModalVisible(false);

      router.replace('/login');

      setTimeout(() => {
        router.push('/login');
      }, 50);
    } catch (err) {
      console.error('Logout failed:', err);
    }
  };

  const cancelLogout = () => {
    setModalVisible(false);
    router.back();
  };

  return (
    <View style={styles.container}>
      <Modal visible={modalVisible} transparent animationType="fade">
        <Pressable style={styles.overlay} onPress={cancelLogout}>
          <View style={styles.modalBox}>
            <LogOut size={46} color="#035284" style={{ marginBottom: 12 }} />
            <Text style={styles.title}>Confirm Logout</Text>
            <Text style={styles.subtitle}>
              Are you sure you want to log out from your account?
            </Text>

            <View style={styles.buttonRow}>
              <TouchableOpacity
                style={[styles.btn, styles.cancelBtn]}
                onPress={cancelLogout}
              >
                <XCircle size={18} color="#035284" />
                <Text style={styles.cancelText}>Cancel</Text>
              </TouchableOpacity>

              <TouchableOpacity
                style={[styles.btn, styles.logoutBtn]}
                onPress={confirmLogout}
              >
                <LogOut size={18} color="#fff" />
                <Text style={styles.logoutText}>Log Out</Text>
              </TouchableOpacity>
            </View>
          </View>
        </Pressable>
      </Modal>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#fff' },

  overlay: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: 'rgba(0,0,0,0.4)',
  },

  modalBox: {
    width: '85%',
    backgroundColor: '#fff',
    borderRadius: 20,
    paddingVertical: 30,
    paddingHorizontal: 25,
    alignItems: 'center',
    elevation: 10,
    shadowColor: '#000',
    shadowOpacity: 0.2,
    shadowOffset: { width: 0, height: 3 },
  },

  title: {
    fontSize: 22,
    fontWeight: '700',
    color: '#035284',
    marginBottom: 8,
  },

  subtitle: {
    fontSize: 15,
    textAlign: 'center',
    color: '#444',
    marginBottom: 25,
    lineHeight: 20,
  },

  buttonRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    width: '100%',
  },

  btn: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    flex: 1,
    paddingVertical: 12,
    borderRadius: 8,
    marginHorizontal: 6,
    gap: 6,
  },

  cancelBtn: {
    backgroundColor: '#E6EEF3',
  },
  logoutBtn: {
    backgroundColor: '#035284',
  },

  cancelText: {
    fontSize: 16,
    fontWeight: '600',
    color: '#035284',
  },
  logoutText: {
    fontSize: 16,
    fontWeight: '600',
    color: '#fff',
  },
});
