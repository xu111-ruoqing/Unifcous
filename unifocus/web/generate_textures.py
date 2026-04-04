import os
import numpy as np
from PIL import Image, ImageFilter

out_dir = "public/textures"
os.makedirs(out_dir, exist_ok=True)

width, height = 1024, 512
x = np.linspace(0, 1, width)
y = np.linspace(0, 1, height)
X, Y = np.meshgrid(x, y)

theta = X * 2 * np.pi
phi = Y * np.pi
ny = np.cos(phi)

def fbm(x_coords, y_coords, octaves=4, scale=1.0):
    val = np.zeros_like(x_coords)
    for i in range(octaves):
        freq = (2 ** i) * scale
        amp = 0.5 ** i
        px, py = np.random.rand() * 100, np.random.rand() * 100
        val += amp * np.sin(x_coords * freq + px) * np.cos(y_coords * freq + py)
    return val

print("1. 生成吸光暗物质核心 (Dark Core)...")
# 保留原本 Aleph-1 的风暴走向，但将高能区转为纯黑
storm = fbm(theta, ny * 3, octaves=4, scale=3)
intensity = (np.sin(storm * 4) + 1) / 2
# 反转：风暴中心(1.0)变成纯黑，外围透出幽暗的紫色
r_core = (1 - intensity) * 50
g_core = (1 - intensity) * 20
b_core = (1 - intensity) * 80
# 彻底吞噬光的黑斑
black_holes = intensity > 0.75
r_core[black_holes] = 0
g_core[black_holes] = 0
b_core[black_holes] = 0
core = np.clip(np.stack([r_core, g_core, b_core], axis=2), 0, 255).astype(np.uint8)
Image.fromarray(core).filter(ImageFilter.GaussianBlur(1)).save(os.path.join(out_dir, "core.jpg"))

print("2. 生成引力透镜吸积盘 (Accretion Disk)...")
# 生成一个 1024x1024 的带有 Alpha 透明通道的白炽光环
X_disk, Y_disk = np.meshgrid(np.linspace(-1, 1, 1024), np.linspace(-1, 1, 1024))
R_disk = np.sqrt(X_disk**2 + Y_disk**2)
# 光环厚度与衰减
ring_alpha = np.exp(-((R_disk - 0.55) ** 2) * 45)
# 边缘加入吸积物质的噪点破碎感
angle_disk = np.arctan2(Y_disk, X_disk)
disk_noise = fbm(angle_disk, R_disk, octaves=3, scale=8)
ring_alpha = ring_alpha * (0.8 + 0.3 * disk_noise)
ring_alpha = np.clip(ring_alpha, 0, 1) * 255

r_disk = np.full_like(ring_alpha, 255)
g_disk = np.full_like(ring_alpha, 240)
b_disk = np.full_like(ring_alpha, 255)
disk = np.stack([r_disk, g_disk, b_disk, ring_alpha], axis=2).astype(np.uint8)
Image.fromarray(disk, 'RGBA').save(os.path.join(out_dir, "disk.png"))

print("3. 生成真实的平滑黄褐色土星 (Saturn)...")
# 极低频率的平滑扰动，去除斑马线条纹
distorted_ny_sat = ny + np.sin(theta * 2) * 0.05
bands = np.sin(distorted_ny_sat * 8) * 0.5 + 0.5
r_sat = 200 + bands * 30
g_sat = 170 + bands * 40
b_sat = 120 + bands * 50
saturn = np.clip(np.stack([r_sat, g_sat, b_sat], axis=2), 0, 255).astype(np.uint8)
Image.fromarray(saturn).filter(ImageFilter.GaussianBlur(3)).save(os.path.join(out_dir, "saturn.jpg"))

print("4. 生成深邃无痕的海王星 (Neptune)...")
# 几乎没有条纹，只有大面积的深渊蓝调和隐约的风暴眼
storms = fbm(theta, ny, octaves=3, scale=2)
storms = (storms - storms.min()) / (storms.max() - storms.min())
r_nep = 20 + storms * 15
g_nep = 50 + storms * 25
b_nep = 130 + storms * 40
nep = np.clip(np.stack([r_nep, g_nep, b_nep], axis=2), 0, 255).astype(np.uint8)
Image.fromarray(nep).filter(ImageFilter.GaussianBlur(4)).save(os.path.join(out_dir, "neptune.jpg"))

print("5. 生成紫黑相间的阿尔卑斯大糖果 (Candy)...")
# 放宽间距，强烈的紫黑对比
stripe = np.sin(theta * 2 + ny * 3)
candy_mask = (stripe > 0).astype(float)
r_can = candy_mask * 110 + 15
g_can = candy_mask * 30 + 10
b_can = candy_mask * 200 + 20
candy = np.clip(np.stack([r_can, g_can, b_can], axis=2), 0, 255).astype(np.uint8)
Image.fromarray(candy).filter(ImageFilter.GaussianBlur(1)).save(os.path.join(out_dir, "candy.jpg"))

print("✅ 完美的纹理全部生成！")