import { useEffect, useMemo, useRef, useState } from 'react'
import { categoryApi } from '@/shared/api/v1/category'
import type { Category, CategoryType } from '@/shared/types/category'
import './CategoriesView.css'

export default function CategoriesView() {
  const [types, setTypes] = useState<CategoryType[]>([])
  const [children, setChildren] = useState<Category[]>([])
  const [selectedTypeId, setSelectedTypeId] = useState<number | null>(null)
  const [typeName, setTypeName] = useState('')
  const [categoryName, setCategoryName] = useState('')
  const [notice, setNotice] = useState('')
  const [submittingType, setSubmittingType] = useState(false)
  const [submittingCategory, setSubmittingCategory] = useState(false)
  // ref 作为同步锁：React 的 state 更新是异步的，快速连点时两次点击可能都看到 submitting=false。
  const typeSubmitLock = useRef(false)
  const categorySubmitLock = useRef(false)
  const selectedType = useMemo(() => types.find((item) => item.id === selectedTypeId) ?? null, [types, selectedTypeId])
  const showNotice = (message: string) => { setNotice(message); window.setTimeout(() => setNotice(''), 4000) }
  const errorMessage = (error: unknown) => {
    const message = error instanceof Error ? error.message : '操作失败'
    return /duplicate|exists|1062|唯一/i.test(message) ? '该分类已存在，请勿重复添加' : message
  }
  const loadTypes = async () => { const result = await categoryApi.listTypes(); const nextTypes = result.types ?? []; setTypes(nextTypes); setSelectedTypeId((current) => current && nextTypes.some((item) => item.id === current) ? current : nextTypes[0]?.id ?? null) }
  const loadChildren = async (parentID: number) => { const result = await categoryApi.listCategories({ parentID }); setChildren(result.categories ?? []) }
  useEffect(() => { void loadTypes() }, [])
  useEffect(() => { if (selectedTypeId !== null) void loadChildren(selectedTypeId); else setChildren([]) }, [selectedTypeId])
  const addType = async () => {
    const name = typeName.trim()
    if (!name || typeSubmitLock.current) return
    typeSubmitLock.current = true
    setSubmittingType(true)
    try {
      await categoryApi.createType({ name })
      setTypeName('')
      showNotice('一级分类已添加')
    } catch (error) {
      showNotice(errorMessage(error))
    } finally {
      // 无论请求成功还是失败都重新拉列表：即使创建响应丢失或命中缓存，也能同步到最新数据。
      await loadTypes().catch(() => undefined)
      typeSubmitLock.current = false
      setSubmittingType(false)
    }
  }
  const removeType = async (id: number, name: string) => { if (!window.confirm(`确定要删除一级分类「${name}」吗？\n\n删除后该分类下的二级分类也会一起移除。`)) return; await categoryApi.deleteType({ id }); showNotice('一级分类已删除'); await loadTypes() }
  const addCategory = async () => {
    if (selectedTypeId === null || categorySubmitLock.current) return
    const name = categoryName.trim()
    if (!name) return
    const parentID = selectedTypeId
    categorySubmitLock.current = true
    setSubmittingCategory(true)
    try {
      await categoryApi.createCategory({ parentID, name })
      setCategoryName('')
      showNotice('二级分类已添加')
    } catch (error) {
      showNotice(errorMessage(error))
    } finally {
      await loadChildren(parentID).catch(() => undefined)
      categorySubmitLock.current = false
      setSubmittingCategory(false)
    }
  }
  const removeCategory = async (id: number, name: string) => { if (!window.confirm(`确定要删除二级分类「${name}」吗？`)) return; await categoryApi.deleteCategory({ id }); showNotice('二级分类已删除'); if (selectedTypeId !== null) await loadChildren(selectedTypeId) }

  return <section className="categories-view">
    {notice && <div className="category-toast" role="status"><span className="toast-check">✓</span>{notice}</div>}
    <div className="taxonomy-toolbar"><div><b>一级分类</b><span>{types.length} 个内容方向</span></div><div className="taxonomy-add"><input aria-label="新增一级分类" value={typeName} disabled={submittingType} onChange={(event) => setTypeName(event.target.value)} onKeyDown={(event) => { if (event.key === 'Enter') void addType() }} placeholder="添加一个内容方向" /><button type="button" disabled={submittingType} onClick={() => void addType()}>{submittingType ? '添加中…' : '新增一级分类'}</button></div></div>
    <div className="type-board">{types.map((item, index) => <article className={`type-card ${selectedTypeId === item.id ? 'active' : ''}`} key={item.id}><button className="type-card-main" type="button" onClick={() => setSelectedTypeId(item.id)}><span className="type-card-number">0{index + 1}</span><span className="type-card-name">{item.name}</span><span className="type-card-meta">{selectedTypeId === item.id ? '正在查看' : '查看子分类'} <span>↗</span></span></button><button className="type-card-delete" type="button" aria-label={`删除${item.name}`} onClick={() => void removeType(item.id, item.name)}>×</button></article>)}{types.length === 0 && <div className="board-empty">还没有一级分类，请从上方开始添加。</div>}</div>
    <section className="subtype-section"><div className="subtype-heading"><div><h3>{selectedType ? selectedType.name : '请选择一个一级分类'}<small className="subtype-total">{selectedType ? children.length : 0} 个子分类</small></h3></div></div>{selectedType ? <><div className="subtype-toolbar"><span>二级分类 <small>用于更细致地组织内容</small></span><div className="subtype-add"><input aria-label="新增二级分类" value={categoryName} disabled={submittingCategory} onChange={(event) => setCategoryName(event.target.value)} onKeyDown={(event) => { if (event.key === 'Enter') void addCategory() }} placeholder="输入二级分类名称" /><button type="button" disabled={submittingCategory} onClick={() => void addCategory()}>{submittingCategory ? '添加中…' : '添加'}</button></div></div><div className="subtype-list">{children.map((item, index) => <article className="type-card subtype-card" key={item.id}><div className="type-card-main"><span className="type-card-number">0{index + 1}</span><span className="type-card-name">{item.name}</span><span className="type-card-meta">{selectedType.name} / {item.name}</span></div><button className="type-card-delete" type="button" aria-label={`删除${item.name}`} onClick={() => void removeCategory(item.id, item.name)}>×</button></article>)}{children.length === 0 && <div className="subtype-empty"><span>＋</span><b>这个方向还没有子分类</b><small>添加后，文章发布时就可以使用它。</small></div>}</div></> : <div className="subtype-empty choose-empty"><span>↗</span><b>从上方选择一个一级分类</b><small>选中后，这里会展开对应的二级分类。</small></div>}</section>
  </section>
}
