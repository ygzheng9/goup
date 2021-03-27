-- 生产订单子件需求 - 备料仓库存
with po as (
	select a.invCode, sum(a.Qty - a.IssQty) qty
	from v_mom_moallocate a
	 inner join v_mom_orderdetail_rpt d on d.ModID = a.ModID
	 inner join v_mom_order_rpt h  on h.MoID = d.MoID
	 inner join inventory i on i.cInvCode = d.invCode
		where 1 =  1
			-- and h.MoCode = '0000010537'
			and d.Status = '3'  -- 审核
	 group by a.invCode
	),
	inv as (
	select  a.cInvCode invCode, SUM(a.iQuantity) qty
		from currentstock a
	 where a.cWhCode = '07'
		 and a.iQuantity <> 0
	 group by a.cInvCode
	),
	cmb as (
	select a.invCode, (a.qty) qty
		from po a
	union all
	select b.invCode, (-1 * b.qty) qty
		from inv b
	),
	cmb2 as (
	select invCode, SUM(qty) diff
		from cmb
	group by InvCode
		having SUM(qty) <> 0
	)
	select a.InvCode invCode, i.cInvName invName, a.diff
		from cmb2 a
	 inner join Inventory i on i.cInvCode = a.InvCode
	 