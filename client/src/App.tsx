import './App.css'
import {Box, Container, CircularProgress, Divider, Typography, Stack } from "@mui/material";
import {useQuery} from "@tanstack/react-query";


export type Rate = {
    cbn: {
        USD: number,
        NGN: number
    }
    coinmarketcap: {
        USD: number,
        NGN: number
    }
}

function App() {
  const { data: rate, isLoading } = useQuery<Rate>({
      queryKey: ["rate"],
      queryFn: async () => {
        try {
            const res = await fetch("/api/myrate")
            const data = await res.json()

            if(!res.ok){
                throw new Error(data.error || "Something went wrong")
            }
            return data || {}
        } catch (error){
         // console.log(error)
        }
        return null
    }, refetchInterval: 60 * 60 * 1000 // refetch every hour
  })


  return (
      <Container sx={{ paddingX: { xs: 2, sm: 5, md: 10 } }}>
          <Stack
              direction="row"
              spacing={2}
              alignItems="center"
              justifyContent="center"
              sx={{ flexWrap: 'wrap' }}
          >
              <Typography variant="h4" align="center">
                  CBN
              </Typography>
              <Divider
                  orientation="vertical"
                  flexItem
                  sx={{ display: { xs: 'block' }, borderColor: 'red', borderWidth: 2 }}
              />

              <Typography variant="h4" align="center">
                  CoinMarketCap
              </Typography>
          </Stack>
          <Stack   direction="row"
                   spacing={2}
                   alignItems="center"
                   justifyContent="center"
                   sx={{ flexWrap: 'wrap' }}>
          {
             isLoading ? <CircularProgress color="success" /> :
                 rate === undefined || rate === null ?   <Typography variant="h4" align="center" color="red">
                     Something with wrong, please refresh page.
                 </Typography> : <>
                     <Stack
                         direction="row"
                         spacing={2}
                         alignItems="center"
                         justifyContent="space-between"
                         sx={{ flexWrap: 'wrap' }}
                     >
                         <Box>
                             <Typography variant="subtitle1" component="span" fontWeight="bold">
                                 USD
                             </Typography>
                             <Typography variant="h6" component="span" sx={{ color: 'green', ml: 1 }}>
                                 ${rate.cbn.USD}
                             </Typography>
                         </Box>
                         <Typography variant="h6" component="span">
                             {"<>"}
                         </Typography>
                         <Box>
                             <Typography variant="subtitle1" component="span" fontWeight="bold">
                                 NGN
                             </Typography>
                             <Typography variant="h6" component="span" sx={{ color: 'red', ml: 1 }}>
                                 N{rate.cbn.NGN}
                             </Typography>
                         </Box>
                     </Stack>
                     <Divider
                         orientation="vertical"
                         sx={{ borderColor: 'red', borderWidth: 2, my: 2 }}
                     />
                     <Stack
                         direction="row"
                         spacing={2}
                         alignItems="center"
                         justifyContent="space-between"
                         sx={{ flexWrap: 'wrap' }}
                     >
                         <Box>
                             <Typography variant="subtitle1" component="span" fontWeight="bold">
                                 USD
                             </Typography>
                             <Typography variant="h6" component="span" sx={{ color: 'green', ml: 1 }}>
                                 ${rate.coinmarketcap.USD}
                             </Typography>
                         </Box>
                         <Typography variant="h6" component="span">
                             {"<>"}
                         </Typography>
                         <Box>
                             <Typography variant="subtitle1" component="span" fontWeight="bold">
                                 NGN
                             </Typography>
                             <Typography variant="h6" component="span" sx={{ color: 'red', ml: 1 }}>
                                 N{rate.coinmarketcap.NGN}
                             </Typography>
                         </Box>
                     </Stack>
                 </>
          }
          </Stack>

      </Container>
  )
}

export default App
