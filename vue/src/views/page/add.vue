<template>
	<div>
		<Card dis-hover>
			<Row :gutter="2">
				<Col span="20">
				<Form ref="dataForm" :rules="dataRules" :model="dataForm">
					<FormItem style="margin-bottom:15px" prop="title">
						<Input v-model="dataForm.title" placeholder="请输入标题"></Input>
					</FormItem>
					<FormItem style="margin-bottom:15px" prop="path">
						<Row>
							<Col span="3" style="min-width:195px">
							{{prefix}}
							</Col>
							<Col span="8"><Input type="text" v-model="dataForm.path" placeholder="请输入访问路径"></Input>
							</Col>
							<Col span="10" style="min-width:250px"> .html&nbsp;
							<Button type="info" @click="emitDraft" :loading="draftLoading">
								<Icon type="ios-trash" size="20" />存草稿
							</Button>&nbsp;
							<Button type="warning" @click="emitFinish" :loading="finishLoading">
								<Icon type="ios-send" size="20" />发 布
							</Button>
							<Button :to="{name:'page-list'}" style="margin-left: 10px">返回列表</Button>
							</Col>
						</Row>
					</FormItem>
				</Form>
				<!-- markdown -->
				<mavon v-model="dataForm.markdown" @on-change="onChange"></mavon>
				</Col>
				<Col span="4">
				<Card dis-hover :bordered="false">
					<p slot="title">
						<Icon type="ios-settings-outline" />
						设置
					</p>
					<Form ref="optsForm" label-position="top">
						<FormItem label="评论" prop="allow">
							<Checkbox v-model="dataForm.allow">允许评论</Checkbox>
						</FormItem>
						<FormItem label="发布日期">
							<DatePicker v-model="dataForm.created" type="datetime" placeholder="选择发布日期和时间" :clearable="false" :editable="false"></DatePicker>
						</FormItem>
						<br><br>
					</Form>
				</Card>
				</Col>
			</Row>
		</Card>
	</div>
</template>
<script>
import Mavon from "@/components/markdown/Mavon.vue";
import { admPageAdd } from "@/api/page";
export default {
	components: {
		Mavon
	},
	data() {
		return {
			prefix: process.env.VUE_APP_SRV + "/page/",
			draftLoading: false,
			finishLoading: false,
			dataForm: {
				id: 0,
				title: "",
				path: "",
				kind: 2,
				status: 0,
				summary: "",
				cate_id: 0,
				allow: false,
				created: new Date(),
				richtext: "",
				markdown: ""
			},
			dataRules: {
				title: [{ required: true, message: " ", trigger: "blur" }],
				path: [{ required: true, message: " ", trigger: "blur" }]
			}
		};
	},
	methods: {
		onChange(html) {
			console.log(html);
			// 保存 html
			this.dataForm.richtext = html;
		},
		// 存草稿
		emitDraft() {
			this.$refs.dataForm.validate((valid) => {
				if (valid) {
					if (this.dataForm.richtext == "") {
						this.$Message.warning("请填写内容");
						return;
					}
					console.log(this.dataForm);
					this.dataForm.status = 1; //草稿
					this.draftLoading = true;
					admPageAdd(this.dataForm).then((resp) => {
						this.draftLoading = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: "保存成功",
								onClose: () => {}
							});
						} else {
							this.$Message.error({ content: resp.msg, duration: 3 });
						}
					});
				}
			});
		},
		// 发布
		emitFinish() {
			this.$refs.dataForm.validate((valid) => {
				if (valid) {
					if (this.dataForm.richtext == "") {
						this.$Message.warning("请填写内容");
						return;
					}
					console.log(this.dataForm);
					this.dataForm.status = 2; //发布
					this.finishLoading = true;
					admPageAdd(this.dataForm).then((resp) => {
						this.finishLoading = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: "保存成功",
								onClose: () => {}
							});
						} else {
							this.$Message.error({ content: resp.msg, duration: 3 });
						}
					});
				}
			});
		}
	}
};
</script>
