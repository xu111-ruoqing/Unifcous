'use client';

import { useEffect, useState } from 'react';
import {
  Table, Button, Modal, Form, Input, Select, Space, Popconfirm,
  Typography, Tag, message, Tooltip,
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, ReloadOutlined, LinkOutlined } from '@ant-design/icons';

const { Title } = Typography;
const { TextArea } = Input;

interface Competition {
  id: number;
  name: string;
  level: string;
  official_url: string;
  typical_time_window: string;
  timeline_hint: string;
  note: string;
}

const LEVEL_OPTIONS = [
  '国家级A*类', '国家级A类', '国家级B类',
  '省级A类', '省级B类', '省级',
  '国家级', '校级',
];

const levelColor: Record<string, string> = {
  '国家级A*类': 'red',
  '国家级A类': 'volcano',
  '国家级B类': 'orange',
  '省级A类': 'blue',
  '省级B类': 'geekblue',
};

async function requestCompetitions(search: string, levelFilter: string) {
  const params = new URLSearchParams();
  if (search) params.set('q', search);
  if (levelFilter) params.set('level', levelFilter);
  params.set('limit', '500');

  const res = await fetch(`/api/v1/competitions?${params}`);
  const json = await res.json();
  return Array.isArray(json) ? json : [];
}

export default function CompetitionsPage() {
  const [data, setData] = useState<Competition[]>([]);
  const [loading, setLoading] = useState(false);
  const [search, setSearch] = useState('');
  const [levelFilter, setLevelFilter] = useState('');
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<Competition | null>(null);
  const [form] = Form.useForm();
  const [submitting, setSubmitting] = useState(false);

  const fetchData = async () => {
    setLoading(true);
    try {
      setData(await requestCompetitions(search, levelFilter));
    } catch {
      message.error('加载失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    let active = true;

    const loadData = async () => {
      setLoading(true);
      try {
        const nextData = await requestCompetitions(search, levelFilter);
        if (active) {
          setData(nextData);
        }
      } catch {
        if (active) {
          message.error('加载失败');
        }
      } finally {
        if (active) {
          setLoading(false);
        }
      }
    };

    void loadData();

    return () => {
      active = false;
    };
  }, [search, levelFilter]);

  const openCreate = () => {
    setEditing(null);
    form.resetFields();
    setModalOpen(true);
  };

  const openEdit = (row: Competition) => {
    setEditing(row);
    form.setFieldsValue(row);
    setModalOpen(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await fetch(`/api/v1/competitions/${id}`, { method: 'DELETE' });
      message.success('已删除');
      fetchData();
    } catch {
      message.error('删除失败');
    }
  };

  const handleSubmit = async () => {
    const values = await form.validateFields();
    setSubmitting(true);
    try {
      const url = editing ? `/api/v1/competitions/${editing.id}` : '/api/v1/competitions';
      const method = editing ? 'PUT' : 'POST';
      const res = await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(values),
      });
      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || '操作失败');
      }
      message.success(editing ? '更新成功' : '创建成功');
      setModalOpen(false);
      fetchData();
    } catch (e: any) {
      message.error(e.message);
    } finally {
      setSubmitting(false);
    }
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 60,
    },
    {
      title: '竞赛名称',
      dataIndex: 'name',
      ellipsis: true,
      render: (text: string, row: Competition) => (
        <Space>
          <span>{text}</span>
          {row.official_url && (
            <Tooltip title={row.official_url}>
              <a href={row.official_url} target="_blank" rel="noreferrer">
                <LinkOutlined />
              </a>
            </Tooltip>
          )}
        </Space>
      ),
    },
    {
      title: '级别',
      dataIndex: 'level',
      width: 120,
      render: (level: string) => (
        <Tag color={levelColor[level] || 'default'}>{level || '-'}</Tag>
      ),
    },
    {
      title: '典型时间',
      dataIndex: 'typical_time_window',
      width: 140,
      ellipsis: true,
    },
    {
      title: '时间线提示',
      dataIndex: 'timeline_hint',
      ellipsis: true,
    },
    {
      title: '操作',
      width: 100,
      render: (_: any, row: Competition) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => openEdit(row)} />
          <Popconfirm title="确认删除？" onConfirm={() => handleDelete(row.id)}>
            <Button size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div style={{ padding: 24 }}>
      <Title level={3}>竞赛信息管理</Title>

      <Space style={{ marginBottom: 16 }} wrap>
        <Input.Search
          placeholder="搜索竞赛名称"
          allowClear
          style={{ width: 240 }}
          onSearch={setSearch}
          onChange={e => !e.target.value && setSearch('')}
        />
        <Select
          placeholder="按级别筛选"
          allowClear
          style={{ width: 160 }}
          onChange={v => setLevelFilter(v || '')}
          options={LEVEL_OPTIONS.map(l => ({ label: l, value: l }))}
        />
        <Button icon={<ReloadOutlined />} onClick={fetchData}>刷新</Button>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>新增竞赛</Button>
        <span style={{ color: '#888' }}>共 {data.length} 条</span>
      </Space>

      <Table
        rowKey="id"
        columns={columns}
        dataSource={data}
        loading={loading}
        size="middle"
        pagination={{ pageSize: 20, showSizeChanger: true }}
        scroll={{ x: 900 }}
      />

      <Modal
        title={editing ? '编辑竞赛' : '新增竞赛'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={() => setModalOpen(false)}
        confirmLoading={submitting}
        width={600}
        okText="保存"
        cancelText="取消"
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="name" label="竞赛名称" rules={[{ required: true, message: '请输入竞赛名称' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="level" label="级别">
            <Select allowClear options={LEVEL_OPTIONS.map(l => ({ label: l, value: l }))} />
          </Form.Item>
          <Form.Item name="official_url" label="官网链接">
            <Input placeholder="https://..." />
          </Form.Item>
          <Form.Item name="typical_time_window" label="典型时间窗口">
            <Input placeholder="如：3–8月" />
          </Form.Item>
          <Form.Item name="timeline_hint" label="时间线提示">
            <Input placeholder="如：报名：3–4月；决赛：8月" />
          </Form.Item>
          <Form.Item name="note" label="备注">
            <TextArea rows={2} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
