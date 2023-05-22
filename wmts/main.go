package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.GET("/wmts", func(context *gin.Context) {
		context.Header("Content-Type", "application/xml")
		context.Writer.WriteString(wmts)
	})
	log.Fatal(engine.Run(":8081"))
}

var wmts = `<Capabilities
	xmlns="http://www.opengis.net/wmts/1.0"
	xmlns:ows="http://www.opengis.net/ows/1.1"
	xmlns:gml="http://www.opengis.net/gml"
	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xmlns:xlink="http://www.w3.org/1999/xlink" xsi:schemaLocation="http://www.opengis.net/wmts/1.0 http://schemas.opengis.net/wmts/1.0.0/wmtsGetCapabilities_response.xsd" version="1.0.0">
	<ows:ServiceIdentification>
		<ows:tile>在线地图服务</ows:tile>
		<ows:Abstract>基于OGC标准的地图服务</ows:Abstract>
		<key_words>
			<ows:Keyword>OGC</ows:Keyword>
		</key_words>
		<service_type codeSpace="wmts"></service_type>
		<ows:ServiceTypeVersion>1.0.0</ows:ServiceTypeVersion>
		<ows:Fees>none</ows:Fees>
		<ows:AccessConstraints>none</ows:AccessConstraints>
	</ows:ServiceIdentification>
	<ows:ServiceProvider>
		<ows:ProviderName>航天宏图信息技术股份有限公司</ows:ProviderName>
		<ows:ProviderSite>https://piesat.cn/</ows:ProviderSite>
		<ows:ServiceContact>
			<ows:IndividualName>航天宏图</ows:IndividualName>
			<ows:PositionName>PIE-Engine</ows:PositionName>
			<ows:ContactInfo>
				<ows:Phone>
					<ows:Voice>010-82556203</ows:Voice>
					<ows:Facsimile>010-82556203</ows:Facsimile>
				</ows:Phone>
				<ows:Address>
					<ows:DeliveryPoint>北京市海淀区杏石口路益园文化创意产业基地A区1号楼5层</ows:DeliveryPoint>
					<ows:City>北京市</ows:City>
					<ows:AdministrativeArea>北京市</ows:AdministrativeArea>
					<ows:Country>中国</ows:Country>
					<ows:PostalCode>101399</ows:PostalCode>
					<ows:ElectronicMailAddress>piesat.cn</ows:ElectronicMailAddress>
				</ows:Address>
				<ows:OnlineResource xlink:type="simple" xlink:href="https://piesat.cn/"></ows:OnlineResource>
			</ows:ContactInfo>
		</ows:ServiceContact>
	</ows:ServiceProvider>
	<Contents>
		<Layer>
			<ows:Title>webp_demp_2</ows:Title>
			<ows:Identifier>webp_demp_2</ows:Identifier>
			<ows:BoundingBox crs="urn:ogc:def:crs:EPSG:6.18.3:3857">
				<ows:LowerCorner>1.1550290966981048E7 4229341.930781716</ows:LowerCorner>
				<ows:UpperCorner>1.148138219838271E7 4292440.112952556</ows:UpperCorner>
			</ows:BoundingBox>
			<Style xlink:isDefault="true">
				<ows:Identifier>default</ows:Identifier>
			</Style>
			<Format>image/tif&amp;format=jpeg&amp;colorRamp=</Format>
			<TileMatrixSetLink>
				<TileMatrixSet>z0to28</TileMatrixSet>
			</TileMatrixSetLink>
			<ResourceURL format="image/tif&amp;format=jpeg&amp;colorRamp=" resourceType="tile" template="http://127.0.0.1:8080/tile-server/v1?layer=wmts_tiles_web_3857&amp;tilematrix={TileMatrix}&amp;Tilecol={TileCol}&amp;TileRow={TileRow}"></ResourceURL>
		</Layer>
		<TileMatrixSet>
			<ows:Identifier>z0to28</ows:Identifier>
			<ows:SupportedCRS>urn:ogc:def:crs:EPSG:6.18.3:3857</ows:SupportedCRS>
			<TileMatrix>
				<ows:Identifier>0</ows:Identifier>
				<ScaleDenominator>559082264.028717800000</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>1</MatrixWidth>
				<MatrixHeight>1</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>1</ows:Identifier>
				<ScaleDenominator>279541132.014358900000</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>2</MatrixWidth>
				<MatrixHeight>2</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>2</ows:Identifier>
				<ScaleDenominator>139770566.007179440000</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>4</MatrixWidth>
				<MatrixHeight>4</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>3</ows:Identifier>
				<ScaleDenominator>69885283.003589720000</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>8</MatrixWidth>
				<MatrixHeight>8</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>4</ows:Identifier>
				<ScaleDenominator>34942641.501794860000</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>16</MatrixWidth>
				<MatrixHeight>16</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>5</ows:Identifier>
				<ScaleDenominator>17471320.750897430000</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>32</MatrixWidth>
				<MatrixHeight>32</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>6</ows:Identifier>
				<ScaleDenominator>8735660.375448715000</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>64</MatrixWidth>
				<MatrixHeight>64</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>7</ows:Identifier>
				<ScaleDenominator>4367830.187724357500</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>128</MatrixWidth>
				<MatrixHeight>128</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>8</ows:Identifier>
				<ScaleDenominator>2183915.093862178700</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>256</MatrixWidth>
				<MatrixHeight>256</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>9</ows:Identifier>
				<ScaleDenominator>1091957.546931089400</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>512</MatrixWidth>
				<MatrixHeight>512</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>10</ows:Identifier>
				<ScaleDenominator>545978.773465544700</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>1024</MatrixWidth>
				<MatrixHeight>1024</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>11</ows:Identifier>
				<ScaleDenominator>272989.386732772340</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>2048</MatrixWidth>
				<MatrixHeight>2048</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>12</ows:Identifier>
				<ScaleDenominator>136494.693366386170</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>4096</MatrixWidth>
				<MatrixHeight>4096</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>13</ows:Identifier>
				<ScaleDenominator>68247.346683193090</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>8192</MatrixWidth>
				<MatrixHeight>8192</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>14</ows:Identifier>
				<ScaleDenominator>34123.673341596540</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>16384</MatrixWidth>
				<MatrixHeight>16384</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>15</ows:Identifier>
				<ScaleDenominator>17061.836670798270</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>32768</MatrixWidth>
				<MatrixHeight>32768</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>16</ows:Identifier>
				<ScaleDenominator>8530.918335399136</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>65536</MatrixWidth>
				<MatrixHeight>65536</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>17</ows:Identifier>
				<ScaleDenominator>4513.99773337655125</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>131072</MatrixWidth>
				<MatrixHeight>131072</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>18</ows:Identifier>
				<ScaleDenominator>2132.729583849784</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>262144</MatrixWidth>
				<MatrixHeight>262144</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>19</ows:Identifier>
				<ScaleDenominator>1066.364791924892</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>524288</MatrixWidth>
				<MatrixHeight>524288</MatrixHeight>
			</TileMatrix>
			<TileMatrix>
				<ows:Identifier>20</ows:Identifier>
				<ScaleDenominator>533.182395962446</ScaleDenominator>
				<TopLeftCorner>-20037508.34278924 20037508.34278924</TopLeftCorner>
				<TileWidth>256</TileWidth>
				<TileHeight>256</TileHeight>
				<MatrixWidth>1048576</MatrixWidth>
				<MatrixHeight>1048576</MatrixHeight>
			</TileMatrix>
		</TileMatrixSet>
	</Contents>
</Capabilities>`
