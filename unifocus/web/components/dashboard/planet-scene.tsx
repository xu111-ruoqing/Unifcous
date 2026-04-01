"use client";

import { useEffect, useRef } from "react";
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
} from "three";
import type { OrbitGroup, OrbitItem, PlanetSizes } from "./planet-data";

// ─── Color palette (dark cosmic mode) ────────────────────────────────────────

const PALETTE = {
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
} as const;

const ORBIT_LAYOUT = {
  scholarships: { center: [-3.6, 2.05] as [number, number], size: [3.5, 1.85] as [number, number] },
  competitions: { center: [-3.2, -2.5] as [number, number], size: [3.35, 1.75] as [number, number] },
  internships: { center: [3.65, 0.8] as [number, number], size: [3.2, 1.72] as [number, number] },
} as const;

// ─── Three.js helpers ─────────────────────────────────────────────────────────

function createGlowTexture(color: Color): CanvasTexture {
  const size = 256;
  const canvas = document.createElement("canvas");
  canvas.width = size;
  canvas.height = size;
  const ctx = canvas.getContext("2d")!;
  const g = ctx.createRadialGradient(size / 2, size / 2, 8, size / 2, size / 2, size / 2);
  g.addColorStop(0, `rgba(${color.r * 255 | 0},${color.g * 255 | 0},${color.b * 255 | 0},1)`);
  g.addColorStop(0.28, `rgba(${color.r * 255 | 0},${color.g * 255 | 0},${color.b * 255 | 0},0.65)`);
  g.addColorStop(1, "rgba(255,255,255,0)");
  ctx.fillStyle = g;
  ctx.fillRect(0, 0, size, size);
  return new CanvasTexture(canvas);
}

function createParticleTexture(color: Color): CanvasTexture {
  const size = 128;
  const canvas = document.createElement("canvas");
  canvas.width = size;
  canvas.height = size;
  const ctx = canvas.getContext("2d")!;
  const g = ctx.createRadialGradient(size / 2, size / 2, 0, size / 2, size / 2, size / 2);
  g.addColorStop(0, `rgba(${color.r * 255 | 0},${color.g * 255 | 0},${color.b * 255 | 0},1)`);
  g.addColorStop(0.45, `rgba(${color.r * 255 | 0},${color.g * 255 | 0},${color.b * 255 | 0},0.7)`);
  g.addColorStop(1, "rgba(255,255,255,0)");
  ctx.fillStyle = g;
  ctx.fillRect(0, 0, size, size);
  return new CanvasTexture(canvas);
}

function createOrbitLine(sizeX: number, sizeY: number, color: number, tilt = 0): LineLoop {
  const curve = new EllipseCurve(0, 0, sizeX, sizeY, 0, Math.PI * 2, false, 0);
  const points = curve.getPoints(240).map((p) => new Vector3(p.x, p.y, 0));
  const geo = new BufferGeometry().setFromPoints(points);
  const mat = new LineBasicMaterial({ color, transparent: true, opacity: 0.42 });
  const loop = new LineLoop(geo, mat);
  loop.rotation.x = tilt;
  return loop;
}

function createPlanet(radius: number, color: number, glowColor: number): Group {
  const group = new Group();
  const glowTex = createGlowTexture(new Color(glowColor));
  const glowMat = new SpriteMaterial({
    map: glowTex,
    color: glowColor,
    transparent: true,
    opacity: 0.72,
    depthWrite: false,
  });
  const glow = new Sprite(glowMat);
  glow.scale.set(radius * 5.6, radius * 5.6, 1);
  group.add(glow);

  const sphere = new Mesh(
    new SphereGeometry(radius, 48, 48),
    new MeshStandardMaterial({
      color,
      emissive: color,
      emissiveIntensity: 0.72,
      roughness: 0.3,
      metalness: 0.08,
    })
  );
  group.add(sphere);
  return group;
}

interface SatelliteEntry {
  mesh: Mesh;
  angle: number;
  speed: number;
  orbitCenter: [number, number];
  orbitSize: [number, number];
  /** 1.0 = inner orbit (near deadline), 1.3 = outer orbit (far) */
  radiusOffset: number;
}

function createSatellites(
  groupKey: string,
  items: OrbitItem[],
  color: number,
  orbitCenter: [number, number],
  orbitSize: [number, number]
): { group: Group; entries: SatelliteEntry[] } {
  const group = new Group();
  const entries: SatelliteEntry[] = [];
  const step = (Math.PI * 2) / Math.max(items.length, 1);
  const angleOffset = groupKey === "internships" ? 0.55 : groupKey === "scholarships" ? 1.35 : 2.2;

  items.forEach((item, i) => {
    // size: weight maps 0→0.07, 1→0.15
    const radius = 0.07 + item.weight * 0.08;
    const mesh = new Mesh(
      new SphereGeometry(radius, 24, 24),
      new MeshStandardMaterial({ color, emissive: color, emissiveIntensity: 0.55 + item.weight * 0.4, roughness: 0.35 })
    );
    mesh.userData = { title: item.title, context: item.context, groupKey };
    group.add(mesh);

    // timeProximity: high → inner orbit (1.0), low → outer orbit (1.3)
    const radiusOffset = 1.3 - item.timeProximity * 0.3;

    entries.push({
      mesh,
      angle: step * i + angleOffset,
      speed: 0.1 + i * 0.012,
      orbitCenter,
      orbitSize,
      radiusOffset,
    });
  });

  return { group, entries };
}

// ─── Component ────────────────────────────────────────────────────────────────

export interface ActiveNode {
  title: string;
  context: string;
  groupKey: string;
}

export interface PlanetSceneProps {
  orbitGroups: OrbitGroup[];
  planetSizes: PlanetSizes;
  onActiveNodeChange: (node: ActiveNode) => void;
}

export default function PlanetScene({ orbitGroups, planetSizes, onActiveNodeChange }: PlanetSceneProps) {
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    // Renderer
    const renderer = new WebGLRenderer({ antialias: true, alpha: true });
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
    renderer.setSize(container.clientWidth, container.clientHeight);
    renderer.setClearColor(PALETTE.background, 0.88);
    container.innerHTML = "";
    container.appendChild(renderer.domElement);

    // Scene
    const scene = new Scene();
    scene.fog = new Fog(PALETTE.fog, 14, 32);

    // Camera
    const camera = new PerspectiveCamera(36, container.clientWidth / container.clientHeight, 0.1, 100);
    camera.position.set(0, 0, 18);

    // Lights
    const keyLight = new PointLight(0xffffff, 11, 45);
    keyLight.position.set(0, 0, 6);
    const fillLight = new PointLight(0x7db8ff, 6, 35);
    fillLight.position.set(-10, 7, 5);
    scene.add(new AmbientLight(PALETTE.ambient, 1.15), keyLight, fillLight);

    const root = new Group();
    root.rotation.x = 0.3;
    scene.add(root);

    // Star particles
    const particleColor = new Color(PALETTE.particle);
    const particleTex = createParticleTexture(particleColor);
    const count = 1400;
    const pos = new Float32Array(count * 3);
    const sizes = new Float32Array(count);
    for (let i = 0; i < count; i++) {
      pos[i * 3] = (Math.random() - 0.5) * 26;
      pos[i * 3 + 1] = (Math.random() - 0.5) * 18;
      pos[i * 3 + 2] = (Math.random() - 0.5) * 12;
      sizes[i] = 0.7 + Math.random() * 1.25;
    }
    const particleGeo = new BufferGeometry();
    particleGeo.setAttribute("position", new BufferAttribute(pos, 3));
    particleGeo.setAttribute("size", new BufferAttribute(sizes, 1));
    const particles = new Points(
      particleGeo,
      new PointsMaterial({
        color: particleColor,
        map: particleTex,
        transparent: true,
        opacity: 0.48,
        size: 0.09,
        sizeAttenuation: true,
        blending: AdditiveBlending,
        depthWrite: false,
      })
    );
    root.add(particles);

    // Central orbit rings
    [
      createOrbitLine(7.5, 5.1, PALETTE.orbit, 0.02),
      createOrbitLine(5.65, 3.75, PALETTE.orbit, 0.02),
      createOrbitLine(3.95, 2.45, PALETTE.orbit, 0.02),
    ].forEach((o) => root.add(o));

    // Per-planet orbit rings
    (Object.keys(ORBIT_LAYOUT) as Array<keyof typeof ORBIT_LAYOUT>).forEach((key) => {
      const l = ORBIT_LAYOUT[key];
      const a = createOrbitLine(l.size[0], l.size[1], PALETTE.orbit, 0.04);
      const b = createOrbitLine(l.size[0] * 0.72, l.size[1] * 0.66, PALETTE.orbit, 0.04);
      a.position.set(l.center[0], l.center[1], 0);
      b.position.set(l.center[0], l.center[1], 0);
      root.add(a, b);
    });

    // Central UniFocus planet
    const uniPlanet = createPlanet(1.48, PALETTE.main, PALETTE.mainGlow);
    root.add(uniPlanet);

    // Category planets with dynamic sizes
    const scholarshipPlanet = createPlanet(planetSizes.scholarships, PALETTE.scholarships, PALETTE.scholarships);
    scholarshipPlanet.position.set(-3.85, 2.25, 0.25);
    root.add(scholarshipPlanet);

    const competitionPlanet = createPlanet(planetSizes.competitions, PALETTE.competitions, PALETTE.competitions);
    competitionPlanet.position.set(-3.45, -2.7, 0.18);
    root.add(competitionPlanet);

    const internshipPlanet = createPlanet(planetSizes.internships, PALETTE.internships, PALETTE.internships);
    internshipPlanet.position.set(3.7, 0.95, 0.22);
    root.add(internshipPlanet);

    // Satellites
    const GROUP_COLOR: Record<string, number> = {
      scholarships: PALETTE.scholarships,
      competitions: PALETTE.competitions,
      internships: PALETTE.internships,
    };

    const allSatellites: SatelliteEntry[] = [];
    orbitGroups.forEach((group) => {
      const layout = ORBIT_LAYOUT[group.key as keyof typeof ORBIT_LAYOUT];
      if (!layout) return;
      const colorKey = GROUP_COLOR[group.key] ?? PALETTE.internships;
      const { group: satGroup, entries } = createSatellites(
        group.key, group.items, colorKey, layout.center, layout.size
      );
      root.add(satGroup);
      allSatellites.push(...entries);
    });

    // Raycasting
    const raycaster = new Raycaster();
    const pointer = new Vector2(10, 10);
    const handleMove = (e: PointerEvent) => {
      const rect = renderer.domElement.getBoundingClientRect();
      pointer.x = ((e.clientX - rect.left) / rect.width) * 2 - 1;
      pointer.y = -((e.clientY - rect.top) / rect.height) * 2 + 1;
    };
    const handleLeave = () => {
      pointer.set(10, 10);
      renderer.domElement.style.cursor = "default";
    };
    renderer.domElement.addEventListener("pointermove", handleMove);
    renderer.domElement.addEventListener("pointerleave", handleLeave);

    // Pre-compute mesh list once so raycaster doesn't allocate a new array every frame
    const satelliteMeshes = allSatellites.map((s) => s.mesh);

    // Animate
    const clock = new Clock();
    let frameId = 0;
    let lastHoveredTitle = "";

    const animate = () => {
      const t = clock.getElapsedTime();

      uniPlanet.rotation.y += 0.006;
      scholarshipPlanet.rotation.y += 0.0042;
      competitionPlanet.rotation.y += 0.0035;
      internshipPlanet.rotation.y += 0.0045;

      uniPlanet.position.y = Math.sin(t * 0.55) * 0.12;
      scholarshipPlanet.position.y = 2.2 + Math.sin(t * 0.72) * 0.16;
      competitionPlanet.position.y = -2.68 + Math.cos(t * 0.68) * 0.15;
      internshipPlanet.position.y = 0.92 + Math.sin(t * 0.63) * 0.17;
      particles.rotation.z += 0.00045;
      root.rotation.z = Math.sin(t * 0.1) * 0.035;

      allSatellites.forEach((sat, idx) => {
        const angle = sat.angle + t * sat.speed;
        const w = sat.orbitSize[0] * sat.radiusOffset;
        const h = sat.orbitSize[1] * sat.radiusOffset;
        sat.mesh.position.set(
          sat.orbitCenter[0] + Math.cos(angle) * w,
          sat.orbitCenter[1] + Math.sin(angle) * h,
          Math.sin(t * 0.9 + idx) * 0.36
        );
      });

      raycaster.setFromCamera(pointer, camera);
      const hits = raycaster.intersectObjects(satelliteMeshes, false);
      if (hits.length) {
        const node = hits[0].object.userData as ActiveNode;
        if (node?.title && node.title !== lastHoveredTitle) {
          lastHoveredTitle = node.title;
          onActiveNodeChange(node);
        }
        renderer.domElement.style.cursor = "pointer";
      } else {
        renderer.domElement.style.cursor = "default";
      }

      renderer.render(scene, camera);
      frameId = requestAnimationFrame(animate);
    };
    animate();

    // Resize
    const ro = new ResizeObserver(() => {
      camera.aspect = container.clientWidth / container.clientHeight;
      camera.updateProjectionMatrix();
      renderer.setSize(container.clientWidth, container.clientHeight);
    });
    ro.observe(container);

    return () => {
      cancelAnimationFrame(frameId);
      ro.disconnect();
      renderer.domElement.removeEventListener("pointermove", handleMove);
      renderer.domElement.removeEventListener("pointerleave", handleLeave);
      scene.traverse((child) => {
        if ("geometry" in child && child.geometry) (child as Mesh).geometry.dispose();
        if ("material" in child && child.material) {
          const mats = Array.isArray((child as Mesh).material) ? (child as Mesh).material as MeshStandardMaterial[] : [(child as Mesh).material as MeshStandardMaterial];
          mats.forEach((m) => { m.map?.dispose(); m.dispose(); });
        }
      });
      renderer.dispose();
    };
  }, [orbitGroups, planetSizes, onActiveNodeChange]);

  return <div ref={containerRef} className="h-full w-full" />;
}
