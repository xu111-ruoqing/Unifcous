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
  DirectionalLight,
  DoubleSide,
  EllipseCurve,
  Fog,
  Group,
  LineBasicMaterial,
  LineLoop,
  Mesh,
  MeshBasicMaterial,
  MeshPhysicalMaterial,
  MeshStandardMaterial,
  PerspectiveCamera,
  Points,
  PointsMaterial,
  Raycaster,
  RingGeometry,
  Scene,
  ShaderMaterial,
  SphereGeometry,
  Sprite,
  SpriteMaterial,
  TextureLoader,
  Vector2,
  Vector3,
  WebGLRenderer,
  SRGBColorSpace,
} from "three";
import type { OrbitGroup, OrbitItem, PlanetSizes } from "./planet-data";

// ─── Color palette ────────────────────────────────────────────────────────
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

// ─── Helpers ─────────────────────────────────────────────────────────

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

function createProceduralBumpMap(): CanvasTexture {
  const size = 512;
  const canvas = document.createElement("canvas");
  canvas.width = size;
  canvas.height = size;
  const ctx = canvas.getContext("2d")!;
  ctx.fillStyle = "#888888";
  ctx.fillRect(0, 0, size, size);
  for (let i = 0; i < 1500; i++) {
    const x = Math.random() * size;
    const y = Math.random() * size;
    const r = Math.random() * 6 + 1;
    const c = Math.floor(Math.random() * 255);
    ctx.beginPath();
    ctx.arc(x, y, r, 0, Math.PI * 2);
    ctx.fillStyle = `rgba(${c},${c},${c},0.12)`;
    ctx.fill();
  }
  return new CanvasTexture(canvas);
}

function createOrbitLine(sizeX: number, sizeY: number, color: number, tilt = 0): LineLoop {
  const curve = new EllipseCurve(0, 0, sizeX, sizeY, 0, Math.PI * 2, false, 0);
  const points = curve.getPoints(240).map((p) => new Vector3(p.x, p.y, 0));
  const geo = new BufferGeometry().setFromPoints(points);
  const mat = new LineBasicMaterial({ 
    color, transparent: true, opacity: 0.18, blending: AdditiveBlending, depthWrite: false 
  });
  const loop = new LineLoop(geo, mat);
  loop.rotation.x = tilt;
  return loop;
}

// 核心改动：创建内存中的动态 Canvas 文本贴图（Billboard 标签）
function createTextSprite(message: string, colorHex: number): Sprite {
  const canvas = document.createElement("canvas");
  canvas.width = 256;
  canvas.height = 128;
  const ctx = canvas.getContext("2d")!;
  
  // 排版样式
  ctx.font = "bold 44px 'PingFang SC', 'Microsoft YaHei', sans-serif";
  ctx.textAlign = "center";
  ctx.textBaseline = "middle";
  
  // 绘制描边（模拟发光，增加太空背景下的可读性）
  const hexStr = "#" + colorHex.toString(16).padStart(6, '0');
  ctx.lineWidth = 6;
  ctx.strokeStyle = "rgba(0, 0, 0, 0.85)";
  ctx.strokeText(message, 128, 64);
  
  // 绘制纯白填充和光晕阴影
  ctx.fillStyle = "#ffffff";
  ctx.shadowColor = hexStr;
  ctx.shadowBlur = 10;
  ctx.fillText(message, 128, 64);
  
  const texture = new CanvasTexture(canvas);
  texture.colorSpace = SRGBColorSpace;
  
  // 注意：保持 depthTest: true，以便标签绕到黑洞背面时能被实体物理遮挡
  const spriteMat = new SpriteMaterial({ 
    map: texture, 
    transparent: true, 
    opacity: 0.9, 
    depthWrite: false 
  });
  
  const sprite = new Sprite(spriteMat);
  sprite.scale.set(0.65, 0.325, 1);
  return sprite;
}

function createPlanet(radius: number, color: number, glowColor: number, bumpMap: CanvasTexture, textureUrl?: string, isBlackHole: boolean = false, hasRing: boolean = false): Group {
  const group = new Group();

  if (isBlackHole) {
    const texLoader = new TextureLoader();
    const diskTex = texLoader.load('/textures/disk.png');
    diskTex.colorSpace = SRGBColorSpace;
    
    // 黑洞本体 + 边缘极限光子圈
    const bhShader = new ShaderMaterial({
      uniforms: {
        glowColor: { value: new Color(0xffeedd) }
      },
      vertexShader: `
        varying vec3 vNormal;
        void main() {
          vNormal = normalize(normalMatrix * normal);
          gl_Position = projectionMatrix * modelViewMatrix * vec4(position, 1.0);
        }
      `,
      fragmentShader: `
        uniform vec3 glowColor;
        varying vec3 vNormal;
        void main() {
          float d = max(0.0, dot(normalize(vNormal), vec3(0.0, 0.0, 1.0)));
          float photonRing = smoothstep(0.08, 0.0, d);
          photonRing = pow(photonRing, 2.0);
          gl_FragColor = vec4(glowColor * photonRing * 3.0, 1.0);
        }
      `,
      transparent: true,
      depthWrite: true, // 写入深度以实现对透镜光和卫星文字的物理遮挡
    });
    const bh = new Mesh(new SphereGeometry(radius, 64, 64), bhShader);
    group.add(bh);

    // 引力透镜边缘（面向屏幕拉伸的背部遮挡层）
    const lensedHalo = new Mesh(
      new RingGeometry(radius * 0.76, radius * 1.65, 64),
      new MeshBasicMaterial({
        map: diskTex,
        color: 0xddddff, 
        side: DoubleSide,
        transparent: true,
        opacity: 0.85,
        blending: AdditiveBlending,
        depthWrite: false,
        depthTest: true 
      })
    );
    lensedHalo.position.z = -0.05; 
    lensedHalo.scale.set(1, 1.33, 1); 
    group.add(lensedHalo);

    // 物理吸积盘
    const accretionDisk = new Mesh(
      new RingGeometry(radius * 1.01, radius * 2.3, 64),
      new MeshBasicMaterial({ 
        map: diskTex,
        color: 0xffe6cc,
        side: DoubleSide,
        transparent: true,
        opacity: 0.95,
        blending: AdditiveBlending,
        depthWrite: false,
        depthTest: true
      })
    );
    accretionDisk.rotation.x = Math.PI / 1.85; 
    accretionDisk.name = "accretionDisk"; 
    group.add(accretionDisk);

    // 外围能量散逸
    const outerGlowTex = createGlowTexture(new Color(0x1a0d33));
    const outerGlow = new Sprite(
      new SpriteMaterial({
        map: outerGlowTex, color: 0xffffff, transparent: true, opacity: 0.5, blending: AdditiveBlending, depthWrite: false,
      })
    );
    outerGlow.scale.set(radius * 3.8, radius * 3.8, 1);
    outerGlow.position.z = -0.1;
    group.add(outerGlow);

    return group;
  }
  
  const glowTex = createGlowTexture(new Color(glowColor));
  const glowMat = new SpriteMaterial({
    map: glowTex, color: glowColor, transparent: true, opacity: 0.5, blending: AdditiveBlending, depthWrite: false,
  });
  const glow = new Sprite(glowMat);
  glow.scale.set(radius * 4.2, radius * 4.2, 1);
  group.add(glow);

  const matParams: any = {
    color,
    emissive: color,
    emissiveIntensity: 0.01, 
    roughness: 0.85,          
    metalness: 0.05,
  };

  if (textureUrl) {
    const tex = new TextureLoader().load(textureUrl);
    tex.colorSpace = SRGBColorSpace; 
    matParams.map = tex;
    matParams.color = new Color(0xffffff);
    matParams.emissive = new Color(0x000000); 
  } else {
    matParams.bumpMap = bumpMap;
    matParams.bumpScale = 0.015;
  }

  const sphere = new Mesh(
    new SphereGeometry(radius, 64, 64),
    new MeshPhysicalMaterial(matParams)
  );
  group.add(sphere);

  if (hasRing) {
    const ringTex = new TextureLoader().load('/textures/disk.png');
    ringTex.colorSpace = SRGBColorSpace;
    const ringGeo = new RingGeometry(radius * 1.3, radius * 2.8, 64);
    const ringMat = new MeshBasicMaterial({
      map: ringTex,
      color: color, 
      side: DoubleSide,
      transparent: true,
      opacity: 0.75,
      blending: AdditiveBlending,
      depthWrite: false
    });
    const ring = new Mesh(ringGeo, ringMat);
    ring.rotation.x = Math.PI / 2.2;
    ring.rotation.y = 0.2; 
    group.add(ring);
  }

  return group;
}

interface SatelliteEntry {
  mesh: Mesh;
  currentAngle: number;
  baseSpeed: number;
  orbitCenter: [number, number];
  orbitSize: [number, number];
  radiusOffset: number;
  groupKey: string;
}

function createSatellites(
  groupKey: string, items: OrbitItem[], color: number, orbitCenter: [number, number], orbitSize: [number, number]
): { group: Group; entries: SatelliteEntry[] } {
  const group = new Group();
  const entries: SatelliteEntry[] = [];
  const step = (Math.PI * 2) / Math.max(items.length, 1);
  const angleOffset = groupKey === "internships" ? 0.55 : groupKey === "scholarships" ? 1.35 : 2.2;

  items.forEach((item, i) => {
    const radius = 0.07 + item.weight * 0.08;
    const baseEmissive = 0.4 + item.weight * 0.3; 

    // 卫星球体实体
    const mesh = new Mesh(
      new SphereGeometry(radius, 24, 24),
      new MeshStandardMaterial({ 
        color, emissive: color, emissiveIntensity: baseEmissive, roughness: 0.25, metalness: 0.8
      })
    );
    mesh.userData = {
      title: item.title,
      context: item.context,
      groupKey,
      baseEmissive,
      url: item.url,
      window: item.window,
      score: item.score,
      passed: item.passed,
      gateReason: item.gateReason,
      daysUntilDeadline: item.daysUntilDeadline,
      scoreComponents: item.scoreComponents,
    };

    // 生成专属简称文字广告牌（Billboard Label）
    // 降级保护：如果数据层异常没有抽出 shortName，使用原 title 截断
    const fallbackText = item.shortName || (item.title ? item.title.substring(0, 4) : "未命名");
    const labelSprite = createTextSprite(fallbackText, color);
    
    // 将文字定位在卫星球体的正上方（基于局部坐标系，因此即便球体放大，文字位置也能完美保持相对距离）
    labelSprite.position.set(0, radius + 0.16, 0);
    mesh.add(labelSprite);

    group.add(mesh);

    entries.push({
      mesh,
      currentAngle: step * i + angleOffset,
      baseSpeed: 0.1 + (i % 10) * 0.012,
      orbitCenter,
      orbitSize,
      radiusOffset: 1.3 - item.timeProximity * 0.3,
      groupKey,
    });
  });

  return { group, entries };
}

export interface ActiveNode {
  title: string;
  context: string;
  groupKey: string;
  url?: string;
  window?: string;
  score?: number;
  passed?: boolean;
  gateReason?: string;
  daysUntilDeadline?: number;
  scoreComponents?: OrbitItem["scoreComponents"];
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

    const renderer = new WebGLRenderer({ antialias: true, alpha: true });
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
    renderer.setSize(container.clientWidth, container.clientHeight);
    renderer.setClearColor(PALETTE.background, 0.88);
    container.innerHTML = "";
    container.appendChild(renderer.domElement);

    const scene = new Scene();
    scene.fog = new Fog(PALETTE.fog, 14, 32);

    const camera = new PerspectiveCamera(36, container.clientWidth / container.clientHeight, 0.1, 100);
    camera.position.set(0, 0, 18);

    const ambientLight = new AmbientLight(PALETTE.ambient, 0.3);
    const keyLight = new DirectionalLight(0xffffff, 4.0);
    keyLight.position.set(8, 6, 8);
    const fillLight = new DirectionalLight(0x52d4ff, 1.5);
    fillLight.position.set(-8, -2, 4);
    const rimLight = new DirectionalLight(PALETTE.mainGlow, 5.0);
    rimLight.position.set(0, 8, -10);
    scene.add(ambientLight, keyLight, fillLight, rimLight);

    const root = new Group();
    root.rotation.x = 0.3;
    scene.add(root);

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
        color: particleColor, map: particleTex, transparent: true, opacity: 0.48, size: 0.09, sizeAttenuation: true, blending: AdditiveBlending, depthWrite: false,
      })
    );
    root.add(particles);

    [
      createOrbitLine(7.5, 5.1, PALETTE.orbit, 0.02),
      createOrbitLine(5.65, 3.75, PALETTE.orbit, 0.02),
      createOrbitLine(3.95, 2.45, PALETTE.orbit, 0.02),
    ].forEach((o) => root.add(o));

    (Object.keys(ORBIT_LAYOUT) as Array<keyof typeof ORBIT_LAYOUT>).forEach((key) => {
      const l = ORBIT_LAYOUT[key];
      const a = createOrbitLine(l.size[0], l.size[1], PALETTE.orbit, 0.04);
      const b = createOrbitLine(l.size[0] * 0.72, l.size[1] * 0.66, PALETTE.orbit, 0.04);
      a.position.set(l.center[0], l.center[1], 0);
      b.position.set(l.center[0], l.center[1], 0);
      root.add(a, b);
    });

    const sharedBumpMap = createProceduralBumpMap();

    const uniPlanet = createPlanet(0.85, PALETTE.main, PALETTE.mainGlow, sharedBumpMap, undefined, true, false);
    root.add(uniPlanet);

    const scholarshipPlanet = createPlanet(planetSizes.scholarships, PALETTE.scholarships, PALETTE.scholarships, sharedBumpMap, '/textures/saturn.jpg', false, true);
    scholarshipPlanet.position.set(-3.85, 2.25, 0.25);
    root.add(scholarshipPlanet);

    const competitionPlanet = createPlanet(planetSizes.competitions, PALETTE.competitions, PALETTE.competitions, sharedBumpMap, '/textures/neptune.jpg');
    competitionPlanet.position.set(-3.45, -2.7, 0.18);
    root.add(competitionPlanet);

    const internshipPlanet = createPlanet(planetSizes.internships, PALETTE.internships, PALETTE.internships, sharedBumpMap, '/textures/candy.jpg');
    internshipPlanet.position.set(3.7, 0.95, 0.22);
    root.add(internshipPlanet);

    const GROUP_COLOR: Record<string, number> = {
      scholarships: PALETTE.scholarships, competitions: PALETTE.competitions, internships: PALETTE.internships,
    };

    const allSatellites: SatelliteEntry[] = [];
    orbitGroups.forEach((group) => {
      const layout = ORBIT_LAYOUT[group.key as keyof typeof ORBIT_LAYOUT];
      if (!layout) return;
      const colorKey = GROUP_COLOR[group.key] ?? PALETTE.internships;
      const { group: satGroup, entries } = createSatellites(group.key, group.items, colorKey, layout.center, layout.size);
      root.add(satGroup);
      allSatellites.push(...entries);
    });

    const raycaster = new Raycaster();
    raycaster.params.Points = { threshold: 0.1 };
    const pointer = new Vector2(10, 10);
    let currentHoveredMesh: Mesh | null = null; 

    const handleMove = (e: PointerEvent) => {
      const rect = renderer.domElement.getBoundingClientRect();
      pointer.x = ((e.clientX - rect.left) / rect.width) * 2 - 1;
      pointer.y = -((e.clientY - rect.top) / rect.height) * 2 + 1;
    };
    
    const handleLeave = () => { pointer.set(10, 10); currentHoveredMesh = null; renderer.domElement.style.cursor = "default"; };
    const handleClick = () => { if (currentHoveredMesh) { const node = currentHoveredMesh.userData as ActiveNode; if (node.url) window.open(node.url, "_blank", "noopener,noreferrer"); } };

    renderer.domElement.addEventListener("pointermove", handleMove);
    renderer.domElement.addEventListener("pointerleave", handleLeave);
    renderer.domElement.addEventListener("click", handleClick);

    const satelliteMeshes = allSatellites.map((s) => s.mesh);
    const targetScaleVec = new Vector3();
    const clock = new Clock();
    let frameId = 0;
    let lastHoveredTitle = "";

    const animate = () => {
      const delta = clock.getDelta();
      const t = clock.getElapsedTime();

      uniPlanet.position.y = Math.sin(t * 0.55) * 0.12; 
      
      const accretionDisk = uniPlanet.getObjectByName("accretionDisk");
      if (accretionDisk) {
        accretionDisk.rotation.z -= 0.003; 
      }

      scholarshipPlanet.rotation.y += 0.002;
      competitionPlanet.rotation.y += 0.0015;
      internshipPlanet.rotation.y += 0.0025;

      scholarshipPlanet.position.y = 2.2 + Math.sin(t * 0.72) * 0.16;
      competitionPlanet.position.y = -2.68 + Math.cos(t * 0.68) * 0.15;
      internshipPlanet.position.y = 0.92 + Math.sin(t * 0.63) * 0.17;

      (scholarshipPlanet.children[0] as Sprite).material.opacity = 0.45 + Math.sin(t * 1.2) * 0.1;
      (competitionPlanet.children[0] as Sprite).material.opacity = 0.45 + Math.sin(t * 1.3) * 0.1;
      (internshipPlanet.children[0] as Sprite).material.opacity = 0.45 + Math.sin(t * 1.4) * 0.1;

      particles.rotation.z += 0.00045;
      root.rotation.z = Math.sin(t * 0.1) * 0.035;

      raycaster.setFromCamera(pointer, camera);
      const hits = raycaster.intersectObjects(satelliteMeshes, false);
      
      if (hits.length > 0) {
        currentHoveredMesh = hits[0].object as Mesh;
        const node = currentHoveredMesh.userData as ActiveNode;
        if (node?.title && node.title !== lastHoveredTitle) { lastHoveredTitle = node.title; onActiveNodeChange(node); }
        renderer.domElement.style.cursor = node.url ? "pointer" : "default"; 
      } else {
        currentHoveredMesh = null;
        renderer.domElement.style.cursor = "default";
      }

      const hoveredGroupKey = currentHoveredMesh ? currentHoveredMesh.userData.groupKey : null;

      allSatellites.forEach((sat, idx) => {
        let currentSpeed = sat.baseSpeed;
        if (currentHoveredMesh) {
          if (sat.mesh === currentHoveredMesh) currentSpeed = 0; 
          else if (sat.groupKey === hoveredGroupKey) currentSpeed = sat.baseSpeed * 0.2; 
        }

        sat.currentAngle += currentSpeed * delta;

        const w = sat.orbitSize[0] * sat.radiusOffset;
        const h = sat.orbitSize[1] * sat.radiusOffset;
        let px = sat.orbitCenter[0] + Math.cos(sat.currentAngle) * w;
        let py = sat.orbitCenter[1] + Math.sin(sat.currentAngle) * h;
        let pz = Math.sin(t * 0.9 + idx) * 0.36;

        const cx = uniPlanet.position.x;
        const cy = uniPlanet.position.y;
        const cz = uniPlanet.position.z;
        const dx = px - cx;
        const dy = py - cy;
        const dz = pz - cz;
        const dist = Math.sqrt(dx * dx + dy * dy + dz * dz);

        const REPULSION_RADIUS = 1.6;
        if (dist < REPULSION_RADIUS) {
          const pushFactor = Math.pow((REPULSION_RADIUS - dist) / REPULSION_RADIUS, 1.8) * 1.6;
          px += (dx / dist) * pushFactor;
          py += (dy / dist) * pushFactor;
          pz += (dz / dist) * pushFactor;
        }

        sat.mesh.position.set(px, py, pz);

        const isHovered = sat.mesh === currentHoveredMesh;
        const targetScale = isHovered ? 1.85 : 1.0;
        const targetEmissive = isHovered ? 2.8 : sat.mesh.userData.baseEmissive;

        // 平滑放缩包含文本标签的卫星组
        targetScaleVec.set(targetScale, targetScale, targetScale);
        sat.mesh.scale.lerp(targetScaleVec, 0.15);
        const mat = sat.mesh.material as MeshStandardMaterial;
        mat.emissiveIntensity += (targetEmissive - mat.emissiveIntensity) * 0.15;
      });

      renderer.render(scene, camera);
      frameId = requestAnimationFrame(animate);
    };
    
    animate();

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
      renderer.domElement.removeEventListener("click", handleClick);
      scene.traverse((child) => {
        if ("geometry" in child && child.geometry) (child as Mesh).geometry.dispose();
        if ("material" in child && child.material) {
          const mats = Array.isArray((child as Mesh).material) ? (child as Mesh).material as MeshStandardMaterial[] : [(child as Mesh).material as MeshStandardMaterial];
          mats.forEach((m) => { m.map?.dispose(); m.dispose(); });
        }
      });
      renderer.dispose();
      sharedBumpMap.dispose(); 
    };
  }, [orbitGroups, planetSizes, onActiveNodeChange]);

  return <div ref={containerRef} className="h-full w-full" />;
}
