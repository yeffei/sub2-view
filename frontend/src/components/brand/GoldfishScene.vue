<template>
  <div ref="sceneElement" class="goldfish-scene" :class="{ 'is-dark': dark, 'is-ready': isReady }" aria-hidden="true">
    <canvas ref="canvasElement" class="goldfish-canvas"></canvas>
    <img class="goldfish-fallback" :src="inkGoldfishTexture" alt="" />
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import * as THREE from 'three'
import inkGoldfishTexture from '@/assets/brand/sst-ink-goldfish-v1.png'

const props = defineProps<{
  dark: boolean
}>()

const sceneElement = ref<HTMLDivElement | null>(null)
const canvasElement = ref<HTMLCanvasElement | null>(null)
const isReady = ref(false)

const fishSize = 0.966
const tailStart = 0.545
const tailWidth = 1 - tailStart

let renderer: THREE.WebGLRenderer | null = null
let scene: THREE.Scene | null = null
let camera: THREE.OrthographicCamera | null = null
let swimGroup: THREE.Group | null = null
let fishRig: THREE.Group | null = null
let tailPivot: THREE.Group | null = null
let bodyFish: THREE.Mesh<THREE.PlaneGeometry, THREE.ShaderMaterial> | null = null
let tailFish: THREE.Mesh<THREE.PlaneGeometry, THREE.ShaderMaterial> | null = null
let fishTexture: THREE.Texture | null = null
let bodyGeometry: THREE.PlaneGeometry | null = null
let tailGeometry: THREE.PlaneGeometry | null = null
let resizeObserver: ResizeObserver | null = null
let visibilityObserver: IntersectionObserver | null = null
let motionQuery: MediaQueryList | null = null
let animationFrame = 0
let isVisible = true
let isRunning = false
let motionTime = 4.8
let previousFrameTime = 0

const bodyUniforms = {
  uTime: { value: 0 },
  uTexture: { value: null as THREE.Texture | null },
  uNight: { value: 0 },
  uStatic: { value: 0 },
  uActivity: { value: 0 },
  uTurning: { value: 0 },
}

const tailUniforms = {
  uTime: { value: 0 },
  uTexture: { value: null as THREE.Texture | null },
  uNight: { value: 0 },
  uStatic: { value: 0 },
  uActivity: { value: 0 },
}

const bodyVertexShader = `
  uniform float uTime;
  uniform float uStatic;
  uniform float uActivity;
  uniform float uTurning;
  varying vec2 vUv;

  void main() {
    float time = mix(uTime, 0.0, uStatic);
    vUv = uv;
    vec3 displaced = position;
    float bodyWeight = smoothstep(0.34, 0.48, uv.x) * (1.0 - smoothstep(0.58, 0.7, uv.x));
    float pectoralWeight = 1.0 - smoothstep(0.035, 0.15, length((uv - vec2(0.49, 0.59)) * vec2(1.0, 0.78)));
    float bodyWave = sin(time * (2.1 + uActivity * 0.55) - uv.x * 5.2) * bodyWeight * (0.0035 + uActivity * 0.0075);
    float finFlutter = sin(time * (3.7 + uTurning * 1.1) + 0.56) * pectoralWeight * (0.004 + uTurning * 0.009);

    displaced.y += bodyWave + finFlutter;
    displaced.x += cos(time * 3.7 + 0.56) * pectoralWeight * (0.0015 + uTurning * 0.003);
    gl_Position = projectionMatrix * modelViewMatrix * vec4(displaced, 1.0);
  }
`

const bodyFragmentShader = `
  uniform sampler2D uTexture;
  uniform float uNight;
  varying vec2 vUv;

  void main() {
    vec4 color = texture2D(uTexture, vUv);
    float bodyMask = 1.0 - smoothstep(0.575, 0.64, vUv.x);
    color.a *= bodyMask;
    color.rgb *= mix(vec3(1.0), vec3(0.88, 0.9, 0.88), uNight);
    gl_FragColor = color;
  }
`

const tailVertexShader = `
  uniform float uTime;
  uniform float uStatic;
  uniform float uActivity;
  varying vec2 vUv;

  void main() {
    float time = mix(uTime, 0.0, uStatic);
    vUv = uv;
    vec3 displaced = position;
    float tipWeight = smoothstep(0.08, 1.0, uv.x);
    float tailWave = sin(time * (4.55 + uActivity * 1.25) - uv.x * 5.1 - 0.45) * tipWeight * (0.008 + uActivity * 0.018);
    float tailDrag = cos(time * (4.55 + uActivity * 1.25) - uv.x * 5.1 - 0.45) * tipWeight * (0.002 + uActivity * 0.004);

    displaced.y += tailWave;
    displaced.x += tailDrag;
    gl_Position = projectionMatrix * modelViewMatrix * vec4(displaced, 1.0);
  }
`

const tailFragmentShader = `
  uniform sampler2D uTexture;
  uniform float uNight;
  varying vec2 vUv;

  void main() {
    vec2 sourceUv = vec2(mix(${tailStart.toFixed(3)}, 1.0, vUv.x), vUv.y);
    vec4 color = texture2D(uTexture, sourceUv);
    float rootFeather = smoothstep(${tailStart.toFixed(3)}, 0.59, sourceUv.x);
    color.a *= rootFeather;
    color.rgb *= mix(vec3(1.0), vec3(0.88, 0.9, 0.88), uNight);
    gl_FragColor = color;
  }
`

function shouldAnimate() {
  return !(motionQuery?.matches ?? true)
}

function resizeScene() {
  if (!renderer || !camera || !sceneElement.value) return

  const { width, height } = sceneElement.value.getBoundingClientRect()
  if (width === 0 || height === 0) return

  const aspect = width / height
  camera.left = -aspect / 2
  camera.right = aspect / 2
  camera.top = 0.5
  camera.bottom = -0.5
  camera.updateProjectionMatrix()
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 1.5))
  renderer.setSize(width, height, false)
}

function easeInOut(value: number) {
  return value * value * (3 - 2 * value)
}

function updateSwimState(time: number) {
  if (!swimGroup || !fishRig || !tailPivot) return

  if (!shouldAnimate()) {
    swimGroup.position.set(0, 0, 0)
    fishRig.rotation.set(0, 0, 0)
    tailPivot.rotation.z = 0
    bodyUniforms.uActivity.value = 0
    bodyUniforms.uTurning.value = 0
    tailUniforms.uActivity.value = 0
    return
  }

  const progress = easeInOut(Math.min(time / 12, 1))
  const x = 0.13 - progress * 0.25
  const y = 0.018 + Math.sin(progress * Math.PI * 1.2 - 0.24) * 0.021
  const activity = progress < 1 ? 0.66 + Math.sin(progress * Math.PI) * 0.16 : 0.22

  swimGroup.position.set(x, y, 0)
  fishRig.rotation.y = 0
  fishRig.rotation.z = Math.sin(time * 1.15) * (progress < 1 ? 0.008 : 0.004)
  tailPivot.rotation.z = Math.sin(time * (4.55 + activity * 1.25) - 0.45) * (0.025 + activity * 0.12)
  bodyUniforms.uActivity.value = activity
  bodyUniforms.uTurning.value = 0
  tailUniforms.uActivity.value = activity
}

function render(time: number) {
  if (!renderer || !scene || !camera) return
  updateSwimState(time)
  bodyUniforms.uTime.value = time
  tailUniforms.uTime.value = time
  renderer.render(scene, camera)
}

function startAnimation() {
  if (isRunning || !shouldAnimate() || !isVisible) return

  isRunning = true
  previousFrameTime = 0
  const animate = (now: number) => {
    if (!isRunning) return
    const delta = previousFrameTime === 0 ? 0 : Math.min((now - previousFrameTime) / 1000, 0.05)
    previousFrameTime = now
    motionTime += delta
    render(motionTime)
    animationFrame = requestAnimationFrame(animate)
  }
  animationFrame = requestAnimationFrame(animate)
}

function stopAnimation() {
  isRunning = false
  previousFrameTime = 0
  cancelAnimationFrame(animationFrame)
  animationFrame = 0
}

function disposeScene() {
  stopAnimation()
  resizeObserver?.disconnect()
  resizeObserver = null
  visibilityObserver?.disconnect()
  visibilityObserver = null
  bodyGeometry?.dispose()
  bodyGeometry = null
  tailGeometry?.dispose()
  tailGeometry = null
  bodyFish?.material.dispose()
  bodyFish = null
  tailFish?.material.dispose()
  tailFish = null
  tailPivot = null
  fishRig = null
  swimGroup = null
  fishTexture?.dispose()
  fishTexture = null
  bodyUniforms.uTexture.value = null
  tailUniforms.uTexture.value = null
  renderer?.dispose()
  renderer?.forceContextLoss()
  renderer = null
  scene = null
  camera = null
  isReady.value = false
}

function initializeScene() {
  if (!canvasElement.value || !sceneElement.value || renderer) return

  try {
    renderer = new THREE.WebGLRenderer({ canvas: canvasElement.value, alpha: true, antialias: true })
    renderer.setClearColor(0x000000, 0)
    renderer.outputColorSpace = THREE.SRGBColorSpace
    scene = new THREE.Scene()
    camera = new THREE.OrthographicCamera(-0.75, 0.75, 0.5, -0.5, 0.1, 10)
    camera.position.z = 2

    swimGroup = new THREE.Group()
    fishRig = new THREE.Group()
    tailPivot = new THREE.Group()
    scene.add(swimGroup)
    swimGroup.add(fishRig)
    fishRig.add(tailPivot)

    bodyGeometry = new THREE.PlaneGeometry(fishSize, fishSize, 64, 18)
    tailGeometry = new THREE.PlaneGeometry(fishSize * tailWidth, fishSize, 32, 18)
    const bodyMaterial = new THREE.ShaderMaterial({
      uniforms: bodyUniforms,
      vertexShader: bodyVertexShader,
      fragmentShader: bodyFragmentShader,
      transparent: true,
      depthWrite: false,
      side: THREE.DoubleSide,
    })
    const tailMaterial = new THREE.ShaderMaterial({
      uniforms: tailUniforms,
      vertexShader: tailVertexShader,
      fragmentShader: tailFragmentShader,
      transparent: true,
      depthWrite: false,
      side: THREE.DoubleSide,
    })

    bodyFish = new THREE.Mesh(bodyGeometry, bodyMaterial)
    bodyFish.position.y = -0.015
    bodyFish.renderOrder = 1
    tailFish = new THREE.Mesh(tailGeometry, tailMaterial)
    tailFish.position.set((fishSize * tailWidth) / 2, -0.015, 0)
    tailFish.renderOrder = 0
    fishRig.add(bodyFish)
    tailPivot.position.set(-fishSize / 2 + fishSize * tailStart, 0, 0)
    tailPivot.add(tailFish)
    resizeScene()

    resizeObserver = new ResizeObserver(resizeScene)
    resizeObserver.observe(sceneElement.value)
    visibilityObserver = new IntersectionObserver(([entry]) => {
      isVisible = entry.isIntersecting
      if (isVisible) startAnimation()
      else stopAnimation()
    }, { threshold: 0.05 })
    visibilityObserver.observe(sceneElement.value)

    new THREE.TextureLoader().load(
      inkGoldfishTexture,
      (loadedTexture: THREE.Texture) => {
        if (!bodyFish || !tailFish || !renderer) {
          loadedTexture.dispose()
          return
        }
        loadedTexture.colorSpace = THREE.SRGBColorSpace
        fishTexture = loadedTexture
        bodyUniforms.uTexture.value = loadedTexture
        tailUniforms.uTexture.value = loadedTexture
        isReady.value = true
        render(motionTime)
        startAnimation()
      },
      undefined,
      () => disposeScene(),
    )
  } catch {
    disposeScene()
  }
}

function handleMotionPreferenceChange() {
  const staticMode = shouldAnimate() ? 0 : 1
  bodyUniforms.uStatic.value = staticMode
  tailUniforms.uStatic.value = staticMode
  if (shouldAnimate()) startAnimation()
  else {
    stopAnimation()
    motionTime = 4.8
    render(motionTime)
  }
}

watch(() => props.dark, (dark) => {
  const night = dark ? 1 : 0
  bodyUniforms.uNight.value = night
  tailUniforms.uNight.value = night
})

onMounted(() => {
  motionQuery = window.matchMedia('(max-width: 1023px), (prefers-reduced-motion: reduce)')
  const night = props.dark ? 1 : 0
  bodyUniforms.uNight.value = night
  tailUniforms.uNight.value = night
  handleMotionPreferenceChange()
  motionQuery.addEventListener('change', handleMotionPreferenceChange)
  initializeScene()
})

onBeforeUnmount(() => {
  motionQuery?.removeEventListener('change', handleMotionPreferenceChange)
  motionQuery = null
  disposeScene()
})
</script>

<style scoped>
.goldfish-scene {
  position: relative;
  width: 100%;
  height: 100%;
  overflow: visible;
  filter: drop-shadow(0 17px 26px rgba(72, 46, 27, 0.09));
}

.goldfish-canvas,
.goldfish-fallback {
  position: absolute;
  inset: 0;
  display: block;
  width: 100%;
  height: 100%;
}

.goldfish-canvas {
  opacity: 0;
  transition: opacity 420ms ease;
}

.goldfish-fallback {
  object-fit: contain;
  transform: scale(0.718);
  transition: opacity 420ms ease;
}

.goldfish-scene.is-ready .goldfish-canvas {
  opacity: 1;
}

.goldfish-scene.is-ready .goldfish-fallback {
  opacity: 0;
}

.goldfish-scene.is-dark {
  filter: brightness(0.94) saturate(0.84) drop-shadow(0 20px 30px rgba(0, 0, 0, 0.28));
}
</style>
