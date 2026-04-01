import { useEffect, useMemo, useState } from 'react';
import PlanetScene from './components/PlanetScene';
import { buildDashboardStats, buildOrbitGroups } from './lib/opportunityData';

const API_STATUS_DEFAULT = {
  competitions: 'idle',
  opportunities: 'idle',
};

const MODE_COPY = {
  dashboard: {
    eyebrow: 'Lightweight Dashboard / 机会直看',
    title: '直接看当前机会，不绕星图。',
    description:
      'Dashboard 回到最初的轻量界面逻辑，优先展示机会本身，让用户先看到内容，再决定是否进入更强叙事的星球视图。',
  },
  planet: {
    eyebrow: 'Opportunity Planet / 展示视图',
    title: '把竞赛、荣誉与实践，放进一张统一的机会星图。',
    description:
      'Planet 模式负责介绍与展示。左侧保留最核心的文案框架，右侧把星球主场景放大，让品牌感和空间感成为第一印象。',
  },
};

const PLANET_STORY = [
  {
    title: '统一机会语言',
    description: '把奖学金、竞赛和实践机会放进同一张图里，用户不用在多个系统之间来回切换理解。',
  },
  {
    title: '更适合介绍',
    description: '星球视图优先承担展示与传播任务，让第一次进入页面的人快速知道 UniFocus 在做什么。',
  },
  {
    title: '后续还能继续接数据',
    description: '当前场景已经可以承接竞赛与机会接口，后面只需要继续扩真实数据映射，不必重做视觉框架。',
  },
];

function formatStatus(status) {
  if (status === 'live') {
    return '在线';
  }

  if (status === 'fallback') {
    return '回落';
  }

  if (status === 'loading') {
    return '加载中';
  }

  return '待连接';
}

export default function App() {
  const [mode, setMode] = useState('dashboard');
  const [competitions, setCompetitions] = useState([]);
  const [opportunities, setOpportunities] = useState([]);
  const [apiStatus, setApiStatus] = useState(API_STATUS_DEFAULT);
  const [activeNode, setActiveNode] = useState({
    title: 'UniFocus 主星',
    context: '右侧星球中的节点会在悬浮时同步刷新这里的说明。',
    groupKey: 'main',
  });

  useEffect(() => {
    const controller = new AbortController();

    async function loadData() {
      setApiStatus({
        competitions: 'loading',
        opportunities: 'loading',
      });

      try {
        const competitionResponse = await fetch('/api/v1/competitions?limit=120', {
          signal: controller.signal,
        });
        if (!competitionResponse.ok) {
          throw new Error('competition api unavailable');
        }
        const competitionData = await competitionResponse.json();
        setCompetitions(Array.isArray(competitionData) ? competitionData : []);
        setApiStatus((previous) => ({
          ...previous,
          competitions: 'live',
        }));
      } catch (error) {
        if (!controller.signal.aborted) {
          setApiStatus((previous) => ({
            ...previous,
            competitions: 'fallback',
          }));
        }
      }

      try {
        const opportunityResponse = await fetch('/api/v1/opportunities?limit=18', {
          signal: controller.signal,
        });
        if (!opportunityResponse.ok) {
          throw new Error('opportunity api unavailable');
        }
        const opportunityData = await opportunityResponse.json();
        setOpportunities(Array.isArray(opportunityData?.data) ? opportunityData.data : []);
        setApiStatus((previous) => ({
          ...previous,
          opportunities: 'live',
        }));
      } catch (error) {
        if (!controller.signal.aborted) {
          setApiStatus((previous) => ({
            ...previous,
            opportunities: 'fallback',
          }));
        }
      }
    }

    void loadData();

    return () => controller.abort();
  }, []);

  const orbitGroups = useMemo(
    () => buildOrbitGroups({ competitions, opportunities }),
    [competitions, opportunities],
  );

  const dashboardStats = useMemo(
    () => buildDashboardStats({ competitions, opportunities, orbitGroups }),
    [competitions, opportunities, orbitGroups],
  );

  const currentCopy = MODE_COPY[mode];

  return (
    <div className={`app-shell app-shell--${mode}`}>
      <div className="app-backdrop" />

      <header className="shell-header">
        <div className="brand-lockup">
          <span className="brand-chip">UniFocus / Opportunity Planet</span>
          <div>
            <h1>机会星球</h1>
            <p>轻量 dashboard 与大场景 planet 的双视角前端</p>
          </div>
        </div>

        <div className="header-actions">
          <div className="mode-toggle" role="tablist" aria-label="view mode">
            {['dashboard', 'planet'].map((item) => (
              <button
                key={item}
                className={item === mode ? 'is-active' : ''}
                onClick={() => setMode(item)}
                type="button"
              >
                {item === 'dashboard' ? 'Dashboard' : 'Planet'}
              </button>
            ))}
          </div>

          <div className="status-strip">
            <span>
              竞赛接口
              <strong>{formatStatus(apiStatus.competitions)}</strong>
            </span>
            <span>
              机会接口
              <strong>{formatStatus(apiStatus.opportunities)}</strong>
            </span>
          </div>
        </div>
      </header>

      {mode === 'dashboard' ? (
        <main className="dashboard-view">
          <section className="dashboard-hero">
            <div className="section-intro">
              <span>{currentCopy.eyebrow}</span>
              <h2>{currentCopy.title}</h2>
              <p>{currentCopy.description}</p>
            </div>

            <div className="dashboard-stats">
              {dashboardStats.map((item) => (
                <article key={item.label} className="stat-pill">
                  <span>{item.label}</span>
                  <strong>{item.value}</strong>
                  <p>{item.detail}</p>
                </article>
              ))}
            </div>
          </section>

          <section className="dashboard-feed surface-panel">
            {orbitGroups.map((group) => (
              <section key={group.key} className={`feed-lane feed-lane--${group.tone}`}>
                <header className="lane-header">
                  <div>
                    <span>{group.english}</span>
                    <h3>{group.label}</h3>
                  </div>
                  <p>{group.summary}</p>
                </header>

                <div className="lane-list">
                  {group.items.map((item) => (
                    <article key={`${group.key}-${item.title}`} className="opportunity-row">
                      <div className="opportunity-row__main">
                        <strong>{item.title}</strong>
                        <p>{item.context}</p>
                      </div>
                      <span>{item.window ?? '开放中'}</span>
                    </article>
                  ))}
                </div>
              </section>
            ))}
          </section>
        </main>
      ) : (
        <main className="planet-view">
          <aside className="planet-copy surface-panel">
            <div className="section-intro">
              <span>{currentCopy.eyebrow}</span>
              <h2>{currentCopy.title}</h2>
              <p>{currentCopy.description}</p>
            </div>

            <div className="story-list">
              {PLANET_STORY.map((item) => (
                <article key={item.title} className="story-line">
                  <strong>{item.title}</strong>
                  <p>{item.description}</p>
                </article>
              ))}
            </div>

            <div className="planet-focus">
              <span>当前焦点</span>
              <strong>{activeNode.title}</strong>
              <p>{activeNode.context}</p>
            </div>
          </aside>

          <section className="planet-stage-shell surface-panel">
            <div className="panel-heading">
              <div>
                <span>主视觉场景</span>
                <h3>Opportunity Galaxy</h3>
              </div>
              <p>右侧保留尽可能大的星球界面，优先承担展示、介绍和氛围塑造。</p>
            </div>

            <PlanetScene
              mode="dashboard"
              orbitGroups={orbitGroups}
              onActiveNodeChange={setActiveNode}
            />
          </section>
        </main>
      )}
    </div>
  );
}
