import { useEffect, useRef } from 'react';
import {
  AdditiveBlending,
  AmbientLight,
  BufferAttribute,
  BufferGeometry,
  CanvasTexture,
  Clock,
  Color,
  EllipseCurve,
  Fog,
  Group,
  LineBasicMaterial,
  LineLoop,
  Mesh,
  MeshStandardMaterial,
  PerspectiveCamera,
  PointLight,
  Points,
  PointsMaterial,
  Raycaster,
  Scene,
  SphereGeometry,
  Sprite,
  SpriteMaterial,
  Vector2,
  Vector3,
  WebGLRenderer,
} from 'three';

const MODE_PRESETS = {
  dashboard: {
    background: 0x040b18,
    fog: 0x081528,
    orbit: 0x3d78ff,
    particle: 0x66c9ff,
    main: 0x7557ff,
    mainGlow: 0xb9a9ff,
    scholarships: 0xe2c471,
    competitions: 0x52d4ff,
    internships: 0xb46dff,
    ambient: 0x9bd1ff,
  },
  planet: {
    background: 0xf4f7ff,
    fog: 0xffffff,
    orbit: 0xc3c9ff,
    particle: 0x8aa3ff,
    main: 0x6e49ff,
    mainGlow: 0xcbc4ff,
    scholarships: 0xf1d896,
    competitions: 0x6be7ff,
    internships: 0xc176ff,
    ambient: 0xffffff,
  },
};

const ORBIT_LAYOUT = {
  scholarships: { center: [-3.6, 2.05], size: [3.5, 1.85], radius: 0.92 },
  competitions: { center: [-3.2, -2.5], size: [3.35, 1.75], radius: 0.94 },
  internships: { center: [3.65, 0.8], size: [3.2, 1.72], radius: 1.0 },
};

function createGlowTexture(color) {
  const size = 256;
  const canvas = document.createElement('canvas');
  canvas.width = size;
  canvas.height = size;
  const context = canvas.getContext('2d');
  const gradient = context.createRadialGradient(size / 2, size / 2, 8, size / 2, size / 2, size / 2);
  gradient.addColorStop(0, `rgba(${color.r * 255}, ${color.g * 255}, ${color.b * 255}, 1)`);
  gradient.addColorStop(0.28, `rgba(${color.r * 255}, ${color.g * 255}, ${color.b * 255}, 0.65)`);
  gradient.addColorStop(1, 'rgba(255, 255, 255, 0)');
  context.fillStyle = gradient;
  context.fillRect(0, 0, size, size);
  return new CanvasTexture(canvas);
}

function createParticleTexture(color) {
  const size = 128;
  const canvas = document.createElement('canvas');
  canvas.width = size;
  canvas.height = size;
  const context = canvas.getContext('2d');
  const gradient = context.createRadialGradient(size / 2, size / 2, 0, size / 2, size / 2, size / 2);
  gradient.addColorStop(0, `rgba(${color.r * 255}, ${color.g * 255}, ${color.b * 255}, 1)`);
  gradient.addColorStop(0.45, `rgba(${color.r * 255}, ${color.g * 255}, ${color.b * 255}, 0.7)`);
  gradient.addColorStop(1, 'rgba(255,255,255,0)');
  context.fillStyle = gradient;
  context.fillRect(0, 0, size, size);
  return new CanvasTexture(canvas);
}

function createOrbitLine(sizeX, sizeY, color, tilt = 0) {
  const curve = new EllipseCurve(0, 0, sizeX, sizeY, 0, Math.PI * 2, false, 0);
  const points = curve.getPoints(240).map((point) => new Vector3(point.x, point.y, 0));
  const geometry = new BufferGeometry().setFromPoints(points);
  const material = new LineBasicMaterial({
    color,
    transparent: true,
    opacity: 0.42,
  });
  const orbit = new LineLoop(geometry, material);
  orbit.rotation.x = tilt;
  return orbit;
}

function createPlanet(radius, color, glowColor, mode) {
  const group = new Group();
  const texture = createGlowTexture(new Color(glowColor));
  const glowMaterial = new SpriteMaterial({
    map: texture,
    color: glowColor,
    transparent: true,
    opacity: mode === 'planet' ? 0.98 : 0.68,
    depthWrite: false,
  });
  const glow = new Sprite(glowMaterial);
  glow.scale.set(radius * 5.6, radius * 5.6, 1);
  group.add(glow);

  const geometry = new SphereGeometry(radius, 48, 48);
  const material = new MeshStandardMaterial({
    color,
    emissive: color,
    emissiveIntensity: mode === 'planet' ? 1.18 : 0.66,
    roughness: 0.3,
    metalness: 0.08,
  });
  const sphere = new Mesh(geometry, material);
  group.add(sphere);
  return group;
}

function createSatellites({ groupKey, items, color, orbitCenter, orbitSize }) {
  const satellites = [];
  const satelliteGroup = new Group();
  const step = (Math.PI * 2) / Math.max(items.length, 1);

  items.forEach((item, index) => {
    const angle = step * index + (groupKey === 'internships' ? 0.55 : groupKey === 'scholarships' ? 1.35 : 2.2);
    const radius = index % 2 === 0 ? 0.13 : 0.095;
    const geometry = new SphereGeometry(radius, 24, 24);
    const material = new MeshStandardMaterial({
      color,
      emissive: color,
      emissiveIntensity: 0.75,
      roughness: 0.35,
    });
    const mesh = new Mesh(geometry, material);
    mesh.userData = { ...item, groupKey };
    satelliteGroup.add(mesh);
    satellites.push({
      mesh,
      angle,
      speed: 0.1 + index * 0.012,
      orbitCenter,
      orbitSize,
      radiusOffset: 1 + (index % 3) * 0.075,
    });
  });

  return { satelliteGroup, satellites };
}

export default function PlanetScene({ mode, orbitGroups, onActiveNodeChange }) {
  const containerRef = useRef(null);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) {
      return undefined;
    }

    const preset = MODE_PRESETS[mode] ?? MODE_PRESETS.dashboard;
    const renderer = new WebGLRenderer({
      antialias: true,
      alpha: true,
    });
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
    renderer.setSize(container.clientWidth, container.clientHeight);
    renderer.setClearColor(preset.background, mode === 'planet' ? 0.05 : 0.88);
    container.innerHTML = '';
    container.appendChild(renderer.domElement);

    const scene = new Scene();
    scene.fog = new Fog(preset.fog, 14, 32);

    const camera = new PerspectiveCamera(36, container.clientWidth / container.clientHeight, 0.1, 100);
    camera.position.set(0, 0, 18);

    const ambientLight = new AmbientLight(preset.ambient, mode === 'planet' ? 2.4 : 1.15);
    const pointLight = new PointLight(0xffffff, mode === 'planet' ? 20 : 11, 45);
    pointLight.position.set(0, 0, 6);
    const sideLight = new PointLight(0x7db8ff, mode === 'planet' ? 11 : 6, 35);
    sideLight.position.set(-10, 7, 5);
    scene.add(ambientLight, pointLight, sideLight);

    const root = new Group();
    root.rotation.x = 0.3;
    scene.add(root);

    const particleColor = new Color(preset.particle);
    const particleTexture = createParticleTexture(particleColor);
    const particleCount = 1400;
    const positions = new Float32Array(particleCount * 3);
    const sizes = new Float32Array(particleCount);

    for (let index = 0; index < particleCount; index += 1) {
      const stride = index * 3;
      positions[stride] = (Math.random() - 0.5) * 26;
      positions[stride + 1] = (Math.random() - 0.5) * 18;
      positions[stride + 2] = (Math.random() - 0.5) * 12;
      sizes[index] = 0.7 + Math.random() * 1.25;
    }

    const particleGeometry = new BufferGeometry();
    particleGeometry.setAttribute('position', new BufferAttribute(positions, 3));
    particleGeometry.setAttribute('size', new BufferAttribute(sizes, 1));
    const particles = new Points(
      particleGeometry,
      new PointsMaterial({
        color: particleColor,
        map: particleTexture,
        transparent: true,
        opacity: mode === 'planet' ? 0.72 : 0.48,
        size: mode === 'planet' ? 0.12 : 0.09,
        sizeAttenuation: true,
        blending: AdditiveBlending,
        depthWrite: false,
      }),
    );
    root.add(particles);

    const centralOrbits = [
      createOrbitLine(7.5, 5.1, preset.orbit, 0.02),
      createOrbitLine(5.65, 3.75, preset.orbit, 0.02),
      createOrbitLine(3.95, 2.45, preset.orbit, 0.02),
    ];
    centralOrbits.forEach((orbit) => root.add(orbit));

    Object.entries(ORBIT_LAYOUT).forEach(([groupKey, layout]) => {
      const orbitA = createOrbitLine(layout.size[0], layout.size[1], preset.orbit, 0.04);
      const orbitB = createOrbitLine(layout.size[0] * 0.72, layout.size[1] * 0.66, preset.orbit, 0.04);
      orbitA.position.set(layout.center[0], layout.center[1], 0);
      orbitB.position.set(layout.center[0], layout.center[1], 0);
      root.add(orbitA, orbitB);
    });

    const uniPlanet = createPlanet(1.48, preset.main, preset.mainGlow, mode);
    root.add(uniPlanet);

    const scholarshipPlanet = createPlanet(0.95, preset.scholarships, preset.scholarships, mode);
    scholarshipPlanet.position.set(-3.85, 2.25, 0.25);
    root.add(scholarshipPlanet);

    const competitionPlanet = createPlanet(1.02, preset.competitions, preset.competitions, mode);
    competitionPlanet.position.set(-3.45, -2.7, 0.18);
    root.add(competitionPlanet);

    const internshipPlanet = createPlanet(1.08, preset.internships, preset.internships, mode);
    internshipPlanet.position.set(3.7, 0.95, 0.22);
    root.add(internshipPlanet);

    const satellites = [];
    orbitGroups.forEach((group) => {
      const layout = ORBIT_LAYOUT[group.key];
      if (!layout) {
        return;
      }

      const colorKey = group.key === 'scholarships'
        ? preset.scholarships
        : group.key === 'competitions'
          ? preset.competitions
          : preset.internships;

      const { satelliteGroup, satellites: groupSatellites } = createSatellites({
        groupKey: group.key,
        items: group.items,
        color: colorKey,
        orbitCenter: layout.center,
        orbitSize: layout.size,
      });
      root.add(satelliteGroup);
      satellites.push(...groupSatellites);
    });

    const raycaster = new Raycaster();
    const pointer = new Vector2(10, 10);

    const handlePointerMove = (event) => {
      const rect = renderer.domElement.getBoundingClientRect();
      pointer.x = ((event.clientX - rect.left) / rect.width) * 2 - 1;
      pointer.y = -((event.clientY - rect.top) / rect.height) * 2 + 1;
    };

    const handlePointerLeave = () => {
      pointer.set(10, 10);
      renderer.domElement.style.cursor = 'default';
    };

    renderer.domElement.addEventListener('pointermove', handlePointerMove);
    renderer.domElement.addEventListener('pointerleave', handlePointerLeave);

    let frameId = 0;
    let lastHoveredTitle = '';
    const clock = new Clock();

    const animate = () => {
      const time = clock.getElapsedTime();

      uniPlanet.rotation.y += 0.006;
      scholarshipPlanet.rotation.y += 0.0042;
      competitionPlanet.rotation.y += 0.0035;
      internshipPlanet.rotation.y += 0.0045;

      uniPlanet.position.y = Math.sin(time * 0.55) * 0.12;
      scholarshipPlanet.position.y = 2.2 + Math.sin(time * 0.72) * 0.16;
      competitionPlanet.position.y = -2.68 + Math.cos(time * 0.68) * 0.15;
      internshipPlanet.position.y = 0.92 + Math.sin(time * 0.63) * 0.17;
      particles.rotation.z += mode === 'planet' ? 0.0008 : 0.00045;
      root.rotation.z = Math.sin(time * 0.1) * (mode === 'planet' ? 0.05 : 0.035);

      satellites.forEach((satellite, index) => {
        const angle = satellite.angle + time * satellite.speed;
        const width = satellite.orbitSize[0] * satellite.radiusOffset;
        const height = satellite.orbitSize[1] * satellite.radiusOffset;
        satellite.mesh.position.set(
          satellite.orbitCenter[0] + Math.cos(angle) * width,
          satellite.orbitCenter[1] + Math.sin(angle) * height,
          Math.sin(time * 0.9 + index) * 0.36,
        );
      });

      raycaster.setFromCamera(pointer, camera);
      const intersections = raycaster.intersectObjects(satellites.map((item) => item.mesh), false);
      if (intersections.length) {
        const nextNode = intersections[0].object.userData;
        if (nextNode?.title && nextNode.title !== lastHoveredTitle) {
          lastHoveredTitle = nextNode.title;
          onActiveNodeChange(nextNode);
        }
        renderer.domElement.style.cursor = 'pointer';
      } else {
        renderer.domElement.style.cursor = 'default';
      }

      renderer.render(scene, camera);
      frameId = window.requestAnimationFrame(animate);
    };

    animate();

    const resizeObserver = new ResizeObserver(() => {
      const width = container.clientWidth;
      const height = container.clientHeight;
      camera.aspect = width / height;
      camera.updateProjectionMatrix();
      renderer.setSize(width, height);
    });

    resizeObserver.observe(container);

    return () => {
      window.cancelAnimationFrame(frameId);
      resizeObserver.disconnect();
      renderer.domElement.removeEventListener('pointermove', handlePointerMove);
      renderer.domElement.removeEventListener('pointerleave', handlePointerLeave);
      scene.traverse((child) => {
        if (child.geometry) {
          child.geometry.dispose();
        }
        if (child.material) {
          const materials = Array.isArray(child.material) ? child.material : [child.material];
          materials.forEach((material) => {
            if (material.map) {
              material.map.dispose();
            }
            material.dispose();
          });
        }
      });
      renderer.dispose();
    };
  }, [mode, onActiveNodeChange, orbitGroups]);

  return <div className="planet-canvas" ref={containerRef} />;
}
