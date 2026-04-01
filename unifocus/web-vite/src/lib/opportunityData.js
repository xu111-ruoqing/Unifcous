const REFERENCE_ITEMS = {
  scholarships: [
    { title: '国家奖学金', accent: 'gold', context: '高势能荣誉识别' },
    { title: '一等奖学金', accent: 'gold', context: '高绩点与成果联动' },
    { title: '励志奖学金', accent: 'gold', context: '成长轨迹加权' },
    { title: '国际级', accent: 'gold', context: '高曝光通道' },
    { title: '校级', accent: 'gold', context: '校内认定资源' },
  ],
  competitions: [
    { title: '中国国际大学生创新大赛', accent: 'cyan', context: '互联网+ 主轨道' },
    { title: '挑战杯', accent: 'cyan', context: '学术科技作品' },
    { title: '数学建模', accent: 'cyan', context: '高频能力入口' },
    { title: '美赛 / 国赛', accent: 'cyan', context: '国际竞赛映射' },
    { title: '创新创业', accent: 'cyan', context: '产品实践联动' },
  ],
  internships: [
    { title: '互联网大厂', accent: 'violet', context: '高竞争密度岗位' },
    { title: '产品经理', accent: 'violet', context: '策略与产品双栖' },
    { title: '研究助理', accent: 'violet', context: '科研路径延展' },
    { title: '科研院所', accent: 'violet', context: '学研协同路径' },
    { title: '运营增长', accent: 'violet', context: '业务落地训练' },
  ],
};

const ORBIT_CONFIG = {
  scholarships: {
    key: 'scholarships',
    label: '奖学金',
    english: 'Scholarships',
    tone: 'gold',
    heroLabel: '学业荣誉轨道',
    summary: '以奖学金和荣誉认定为主，强调成绩、持续性与公开证明。',
  },
  competitions: {
    key: 'competitions',
    label: '学科竞赛',
    english: 'Competitions',
    tone: 'cyan',
    heroLabel: '竞赛成长轨道',
    summary: '由竞赛级别、赛道与时间窗口构成，是最可结构化的机会面。',
  },
  internships: {
    key: 'internships',
    label: '实习实践',
    english: 'Internships',
    tone: 'violet',
    heroLabel: '实践跃迁轨道',
    summary: '偏向岗位和研究实践，适合将能力标签映射成职业试炼。',
  },
};

const TIMELINE_FALLBACK = [
  { month: '03-04', focus: '报名开启', track: '竞赛' },
  { month: '05-06', focus: '材料组装', track: '奖学金' },
  { month: '07-08', focus: '暑期实践', track: '实习' },
  { month: '09-10', focus: '校内认定', track: '奖学金' },
];

function inferOpportunityGroup(item) {
  const type = `${item?.type ?? ''} ${item?.title ?? ''}`.toLowerCase();

  if (type.includes('奖') || type.includes('scholar')) {
    return 'scholarships';
  }

  if (type.includes('实习') || type.includes('intern') || type.includes('岗位') || type.includes('研究')) {
    return 'internships';
  }

  return 'competitions';
}

function compactText(value, fallback) {
  return value && `${value}`.trim() ? `${value}`.trim() : fallback;
}

function sampleCompetitionItems(competitions = []) {
  return competitions.slice(0, 6).map((item) => ({
    title: compactText(item.name, '未命名竞赛'),
    accent: 'cyan',
    context: compactText(item.level, '国家级竞赛'),
    window: compactText(item.typical_time_window, '待补充时间窗口'),
  }));
}

export function buildOrbitGroups({ competitions = [], opportunities = [] }) {
  const byGroup = {
    scholarships: [],
    competitions: sampleCompetitionItems(competitions),
    internships: [],
  };

  opportunities.forEach((item) => {
    const group = inferOpportunityGroup(item);
    byGroup[group].push({
      title: compactText(item.title, '未命名机会'),
      accent: ORBIT_CONFIG[group].tone,
      context: compactText(item.organizer ?? item.type, '机会信号待确认'),
      window: compactText(item.deadline ?? item.start_date, '开放中'),
    });
  });

  return Object.values(ORBIT_CONFIG).map((config) => {
    const items = byGroup[config.key].length
      ? byGroup[config.key].slice(0, 6)
      : REFERENCE_ITEMS[config.key];

    return {
      ...config,
      items,
      count: items.length,
    };
  });
}

export function buildDashboardStats({ competitions = [], opportunities = [], orbitGroups = [] }) {
  const competitionHighLevelCount = competitions.filter((item) => `${item.level ?? ''}`.includes('国家')).length;
  const liveOpportunityCount = opportunities.length;
  const totalOrbitNodes = orbitGroups.reduce((sum, group) => sum + group.items.length, 0);

  return [
    {
      label: '竞赛总量',
      value: String(competitions.length || 84),
      detail: competitionHighLevelCount
        ? `${competitionHighLevelCount} 条国家级数据已接入`
        : '使用仓库基线竞赛数据',
    },
    {
      label: '机会节点',
      value: String(totalOrbitNodes || 15),
      detail: liveOpportunityCount
        ? `后端实时机会 ${liveOpportunityCount} 条`
        : '机会接口未回时使用星球基线节点',
    },
    {
      label: '风格模式',
      value: '2',
      detail: 'Dashboard / Planet 双视角切换',
    },
  ];
}

export function buildTimeline(competitions = []) {
  const extracted = competitions
    .filter((item) => item.typical_time_window)
    .slice(0, 4)
    .map((item) => ({
      month: compactText(item.typical_time_window, '全年'),
      focus: compactText(item.name, '重点竞赛'),
      track: compactText(item.level, '竞赛'),
    }));

  return extracted.length ? extracted : TIMELINE_FALLBACK;
}

export function buildSignalSummary({ competitions = [], opportunities = [] }) {
  return [
    {
      title: '数据联通',
      value: competitions.length ? '竞赛接口在线' : '竞赛接口待连接',
      description: competitions.length
        ? `已抓到 ${competitions.length} 条竞赛数据`
        : '页面会自动回落到仓库内置基线',
    },
    {
      title: '机会映射',
      value: opportunities.length ? '机会接口在线' : '机会接口待连接',
      description: opportunities.length
        ? `已抓到 ${opportunities.length} 条机会记录`
        : '先展示参考图中的典型岗位与荣誉轨道',
    },
  ];
}

export function buildLevelHighlights(competitions = []) {
  const counter = competitions.reduce((acc, item) => {
    const key = compactText(item.level, '未分级');
    acc[key] = (acc[key] || 0) + 1;
    return acc;
  }, {});

  const ranked = Object.entries(counter)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 5)
    .map(([level, count]) => ({
      level,
      count,
    }));

  return ranked.length
    ? ranked
    : [
        { level: '国家级', count: 84 },
        { level: '省级', count: 18 },
        { level: '校级', count: 12 },
      ];
}
