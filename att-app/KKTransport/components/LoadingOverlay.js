import LottieView from 'lottie-react-native';
import { StyleSheet, View } from 'react-native';
import loadingAnim from '../assets/loading_animation.json';

export default function LoadingOverlay({ visible = false }) {
  if (!visible) return null;

  return (
    <View style={styles.overlay}>
      <LottieView
        source={loadingAnim}
        autoPlay
        loop
        style={styles.animation}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  overlay: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(255,255,255,0.85)',
    justifyContent: 'center',
    alignItems: 'center',
    zIndex: 999,
  },
  animation: {
    width: 200,
    height: 200,
  },
});
